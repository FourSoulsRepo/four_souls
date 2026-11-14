package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// maxUnpackBytes caps everything one archive unpacks to.
const maxUnpackBytes = 2 << 30

// unpack extracts the archive into dest and returns the app in it:
// Name.app on macOS, Name.exe on Windows, Name elsewhere.
func unpack(archive, dest, name, goos string) (string, error) {
	if err := os.RemoveAll(dest); err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	if err := os.MkdirAll(dest, 0o700); err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	x := &extractor{root: dest, links: map[string]bool{}, left: maxUnpackBytes}
	var err error
	switch lower := strings.ToLower(archive); {
	case strings.HasSuffix(lower, ".zip"):
		err = x.unzip(archive)
	case strings.HasSuffix(lower, ".tar.gz"):
		err = x.untarGz(archive)
	default:
		err = fmt.Errorf("%w: unknown type %s", ErrBadArchive, filepath.Base(archive))
	}
	if err == nil {
		err = x.makeLinks()
	}
	if err != nil {
		return "", err
	}
	return findPayload(dest, name, goos)
}

// findPayload finds the one app the archive holds; README.txt, LICENSE
// and NOTICE next to it are left out.
func findPayload(root, name, goos string) (string, error) {
	want := name
	switch goos {
	case "darwin":
		want = name + ".app"
	case "windows":
		want = name + ".exe"
	}
	var found string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() != want || p == root {
			return nil
		}
		if goos == "darwin" && d.IsDir() || goos != "darwin" && d.Type().IsRegular() {
			if found == "" {
				found = p
			}
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("%w: no %s in it", ErrBadArchive, want)
	}
	return found, nil
}

// extractor writes archive entries under root and refuses any entry that
// would land outside it, also through a symlink in the archive.
type extractor struct {
	root    string
	links   map[string]bool // relative slash paths of symlinks
	pending []link
	left    int64
}

type link struct {
	rel, target string
}

func (x *extractor) unzip(src string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBadArchive, err)
	}
	defer closeQuiet(r)
	for _, f := range r.File {
		if err := x.zipEntry(f); err != nil {
			return err
		}
	}
	return nil
}

func (x *extractor) zipEntry(f *zip.File) error {
	mode := f.Mode()
	if mode&fs.ModeSymlink != 0 {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrBadArchive, err)
		}
		target, err := io.ReadAll(io.LimitReader(rc, 4096))
		closeQuiet(rc)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrBadArchive, err)
		}
		return x.symlink(f.Name, string(target))
	}
	if mode.IsDir() || strings.HasSuffix(f.Name, "/") {
		return x.dir(f.Name)
	}
	if !mode.IsRegular() {
		return fmt.Errorf("%w: %q is not a file", ErrBadArchive, f.Name)
	}
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBadArchive, err)
	}
	defer closeQuiet(rc)
	return x.file(f.Name, mode, rc)
}

func (x *extractor) untarGz(src string) error {
	f, err := os.Open(src) //nolint:gosec // our own download
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	defer closeQuiet(f)
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBadArchive, err)
	}
	defer closeQuiet(gz)
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("%w: %w", ErrBadArchive, err)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = x.dir(hdr.Name)
		case tar.TypeReg:
			err = x.file(hdr.Name, fs.FileMode(hdr.Mode&0o777), tr) //nolint:gosec // masked to permission bits
		case tar.TypeSymlink:
			err = x.symlink(hdr.Name, hdr.Linkname)
		case tar.TypeXGlobalHeader:
			err = nil
		default:
			err = fmt.Errorf("%w: %q has entry type %q", ErrBadArchive, hdr.Name, hdr.Typeflag)
		}
		if err != nil {
			return err
		}
	}
}

// local turns an entry name into a clean relative slash path, or refuses it.
func (x *extractor) local(name string) (string, error) {
	slash := strings.TrimSuffix(strings.ReplaceAll(name, `\`, "/"), "/")
	if slash == "" || !filepath.IsLocal(filepath.FromSlash(slash)) {
		return "", fmt.Errorf("%w: path %q leaves the folder", ErrBadArchive, name)
	}
	rel := path.Clean(slash)
	// Nothing may be written through a symlink the archive made.
	for p := rel; p != "."; p = path.Dir(p) {
		if x.links[p] {
			return "", fmt.Errorf("%w: path %q goes through a link", ErrBadArchive, name)
		}
	}
	return rel, nil
}

func (x *extractor) dir(name string) error {
	rel, err := x.local(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(x.root, filepath.FromSlash(rel)), 0o755); err != nil { //nolint:gosec // app folders need 0755
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (x *extractor) file(name string, mode fs.FileMode, r io.Reader) error {
	rel, err := x.local(name)
	if err != nil {
		return err
	}
	full := filepath.Join(x.root, filepath.FromSlash(rel))
	if err = os.MkdirAll(filepath.Dir(full), 0o755); err != nil { //nolint:gosec // app folders need 0755
		return fmt.Errorf("update: %w", err)
	}
	out, err := os.OpenFile(full, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm()|0o600) //nolint:gosec // checked by local
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	n, err := io.Copy(out, io.LimitReader(r, x.left+1))
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("update: unpack %s: %w", name, err)
	}
	x.left -= n
	if x.left < 0 {
		return fmt.Errorf("%w: unpacks to over %d bytes", ErrTooLarge, int64(maxUnpackBytes))
	}
	return nil
}

// symlink notes a link; makeLinks makes them once every link is known.
func (x *extractor) symlink(name, target string) error {
	rel, err := x.local(name)
	if err != nil {
		return err
	}
	x.links[rel] = true
	x.pending = append(x.pending, link{rel: rel, target: target})
	return nil
}

// makeLinks makes each link whose target stays inside the folder
// without passing through any link, so no chain of links can leave it.
func (x *extractor) makeLinks() error {
	for _, l := range x.pending {
		t := strings.ReplaceAll(l.target, `\`, "/")
		if t == "" || path.IsAbs(t) || filepath.IsAbs(l.target) || filepath.VolumeName(l.target) != "" {
			return fmt.Errorf("%w: link %q points outside", ErrBadArchive, l.rel)
		}
		// Walk the target from the link's folder, one part at a time.
		var at []string
		if d := path.Dir(l.rel); d != "." {
			at = strings.Split(d, "/")
		}
		for _, part := range strings.Split(t, "/") {
			switch part {
			case "", ".":
			case "..":
				if len(at) == 0 {
					return fmt.Errorf("%w: link %q points outside", ErrBadArchive, l.rel)
				}
				at = at[:len(at)-1]
			default:
				at = append(at, part)
				if x.links[strings.Join(at, "/")] {
					return fmt.Errorf("%w: link %q goes through a link", ErrBadArchive, l.rel)
				}
			}
		}
		full := filepath.Join(x.root, filepath.FromSlash(l.rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil { //nolint:gosec // app folders need 0755
			return fmt.Errorf("update: %w", err)
		}
		if err := os.Symlink(filepath.FromSlash(t), full); err != nil {
			return fmt.Errorf("update: %w", err)
		}
	}
	return nil
}

// closeQuiet closes a reader whose data was already read or dropped.
func closeQuiet(c io.Closer) {
	_ = c.Close() //nolint:errcheck // read only; nothing to do
}
