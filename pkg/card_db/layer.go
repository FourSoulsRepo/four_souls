package carddb

import (
	"io/fs"
	"sort"
	"strings"
)

// WithImages lays extra card images over fsys (16.12), e.g. a
// translator's own images folder. extra holds language folders laid out
// like images/: <lang>/<set>/<file>. A file there wins over the same
// path under images/ in fsys; FindImage then falls back to English as
// before. Card data and every other folder, images/assets included,
// come from fsys only.
func WithImages(fsys, extra fs.FS) fs.FS {
	return imageLayer{base: fsys, extra: extra}
}

type imageLayer struct {
	base, extra fs.FS
}

// extraPath gives the path inside extra for name, when name lies in a
// language folder under images/. Paths that could leave the folder are
// refused, as in card data (ADR 004).
func extraPath(name string) (string, bool) {
	rest, ok := strings.CutPrefix(name, ImagesDir+"/")
	if !ok || !safePath(rest) {
		return "", false
	}
	lang, _, _ := strings.Cut(rest, "/")
	return rest, isLang(lang)
}

// Open serves a file from extra when it has one, else from fsys.
func (l imageLayer) Open(name string) (fs.File, error) {
	var dir fs.File // a folder only extra has
	if sub, ok := extraPath(name); ok {
		if f, err := l.extra.Open(sub); err == nil {
			if info, err := f.Stat(); err == nil && !info.IsDir() {
				return f, nil
			}
			dir = f
		}
	}
	f, err := l.base.Open(name)
	if err != nil && dir != nil {
		return dir, nil
	}
	if dir != nil {
		_ = dir.Close() //nolint:errcheck // a read-only folder
	}
	if err != nil {
		return nil, err //nolint:wrapcheck // keep fs errors as they are
	}
	return f, nil
}

// ReadDir lists a folder of both layers, sorted. In images/ it adds
// extra's language folders, so ImageLangs finds them.
func (l imageLayer) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(l.base, name)
	var more []fs.DirEntry
	var moreErr error
	if name == ImagesDir {
		more, moreErr = fs.ReadDir(l.extra, ".")
		langs := more[:0:0]
		for _, e := range more {
			if e.IsDir() && isLang(e.Name()) {
				langs = append(langs, e)
			}
		}
		more = langs
	} else if sub, ok := extraPath(name); ok {
		more, moreErr = fs.ReadDir(l.extra, sub)
	} else {
		moreErr = fs.ErrNotExist
	}
	if err != nil && moreErr != nil {
		return nil, err //nolint:wrapcheck // keep fs errors as they are
	}
	seen := map[string]bool{}
	for _, e := range entries {
		seen[e.Name()] = true
	}
	for _, e := range more {
		if !seen[e.Name()] {
			entries = append(entries, e)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}
