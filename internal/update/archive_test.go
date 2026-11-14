package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// entry is one archive member for the tests.
type entry struct {
	name string
	body string
	link string // symlink target
	dir  bool
}

func makeZip(t *testing.T, entries []entry) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		body := e.body
		switch {
		case e.dir:
			h.SetMode(fs.ModeDir | 0o755)
		case e.link != "":
			h.SetMode(fs.ModeSymlink | 0o777)
			body = e.link
		default:
			h.SetMode(0o755)
		}
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if !e.dir {
			if _, err := w.Write([]byte(body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func makeTarGz(t *testing.T, entries []entry) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		h := &tar.Header{Name: e.name, Mode: 0o755, Typeflag: tar.TypeReg, Size: int64(len(e.body))}
		switch {
		case e.dir:
			h.Typeflag, h.Size = tar.TypeDir, 0
		case e.link != "":
			h.Typeflag, h.Linkname, h.Size = tar.TypeSymlink, e.link, 0
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if h.Size > 0 {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// releaseEntries is an archive as release.yml packs it: one folder with
// the app, README.txt, LICENSE and NOTICE.
func releaseEntries(goos string) []entry {
	const top = "four-souls-1.2.0-linux-amd64/"
	es := []entry{
		{name: top, dir: true},
		{name: top + "README.txt", body: "readme"},
		{name: top + "LICENSE", body: "mit"},
		{name: top + "NOTICE", body: "notice"},
	}
	switch goos {
	case "darwin":
		return append(es,
			entry{name: top + "FourSouls.app/Contents/MacOS/FourSouls", body: "mach-o"},
			entry{name: top + "FourSouls.app/Contents/Info.plist", body: "plist"})
	case "windows":
		return append(es, entry{name: top + "FourSouls.exe", body: "pe"})
	default:
		return append(es, entry{name: top + "FourSouls", body: "elf"})
	}
}

func unpackBytes(t *testing.T, data []byte, ext, goos string) (string, string, error) {
	t.Helper()
	dir := t.TempDir()
	archive := filepath.Join(dir, "a"+ext)
	if err := os.WriteFile(archive, data, 0o600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	got, err := unpack(archive, out, "FourSouls", goos)
	return got, dir, err
}

func TestUnpackFindsPayload(t *testing.T) {
	for _, tc := range []struct {
		goos, ext, want string
	}{
		{"darwin", ".zip", "FourSouls.app"},
		{"windows", ".zip", "FourSouls.exe"},
		{"linux", ".tar.gz", "FourSouls"},
	} {
		t.Run(tc.goos, func(t *testing.T) {
			data := makeZip(t, releaseEntries(tc.goos))
			if tc.ext == ".tar.gz" {
				data = makeTarGz(t, releaseEntries(tc.goos))
			}
			got, _, err := unpackBytes(t, data, tc.ext, tc.goos)
			if err != nil {
				t.Fatal(err)
			}
			if filepath.Base(got) != tc.want {
				t.Errorf("payload = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUnpackNoPayload(t *testing.T) {
	data := makeTarGz(t, []entry{{name: "dir/README.txt", body: "x"}})
	if _, _, err := unpackBytes(t, data, ".tar.gz", "linux"); !errors.Is(err, ErrBadArchive) {
		t.Errorf("err = %v, want ErrBadArchive", err)
	}
}

func TestUnpackRefusesEscape(t *testing.T) {
	app := entry{name: "top/FourSouls", body: "elf"}
	cases := map[string][]entry{
		"dot dot":            {app, {name: "../evil", body: "x"}},
		"deep dot dot":       {app, {name: "top/../../evil", body: "x"}},
		"absolute":           {app, {name: "/tmp/evil", body: "x"}},
		"backslash":          {app, {name: `..\evil`, body: "x"}},
		"link out":           {app, {name: "top/l", link: "../../evil"}},
		"absolute link":      {app, {name: "top/l", link: "/etc"}},
		"write through link": {app, {name: "top/l", link: "."}, {name: "top/l/x", body: "x"}},
		// d/b points up one level; a then climbs through it out of the folder.
		"link chain": {app, {name: "d/a", link: "b/../.."}, {name: "d/b", link: ".."}},
	}
	for name, es := range cases {
		for _, ext := range []string{".zip", ".tar.gz"} {
			t.Run(name+ext, func(t *testing.T) {
				data := makeZip(t, es)
				if ext == ".tar.gz" {
					data = makeTarGz(t, es)
				}
				_, dir, err := unpackBytes(t, data, ext, "linux")
				if !errors.Is(err, ErrBadArchive) {
					t.Fatalf("err = %v, want ErrBadArchive", err)
				}
				if _, err := os.Lstat(filepath.Join(dir, "evil")); err == nil {
					t.Error("a file was written outside the folder")
				}
			})
		}
	}
}

func TestUnpackKeepsInsideLinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need extra rights on Windows")
	}
	es := []entry{
		{name: "top/FourSouls.app/Contents/MacOS/FourSouls", body: "mach-o"},
		{name: "top/FourSouls.app/Contents/Current", link: "MacOS"},
	}
	for _, ext := range []string{".zip", ".tar.gz"} {
		data := makeZip(t, es)
		if ext == ".tar.gz" {
			data = makeTarGz(t, es)
		}
		got, _, err := unpackBytes(t, data, ext, "darwin")
		if err != nil {
			t.Fatalf("%s: %v", ext, err)
		}
		target, err := os.Readlink(filepath.Join(got, "Contents", "Current"))
		if err != nil || target != "MacOS" {
			t.Errorf("%s: link = %q, %v", ext, target, err)
		}
	}
}
