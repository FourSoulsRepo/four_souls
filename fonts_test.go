package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/fstest"
)

func TestFontsFromFiles(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{
		"images/" + titleFont:    {Data: []byte("\x00\x01\x00\x00")},
		"images/" + bodyFont:     {Data: []byte("\x00\x01\x00\x00")},
		"images/" + cyrillicFont: {Data: []byte("\x00\x01\x00\x00")},
	}}
	got := a.Fonts()
	if got.Title != "/cardimg/"+titleFont || got.Body != "/cardimg/"+bodyFont || got.Cyrillic != "/cardimg/"+cyrillicFont {
		t.Fatalf("%+v", got)
	}
	// Without the private images the system font is used.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	if got := a.Fonts(); got != (Fonts{}) {
		t.Fatalf("%+v without images", got)
	}
}

// The fonts are TrueType in the private images submodule. With the
// submodule they are found and served as fonts; without it, nothing is.
func TestFontsServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	fonts := (&App{cardFiles: files}).Fonts()
	for file, url := range map[string]string{titleFont: fonts.Title, bodyFont: fonts.Body, cyrillicFont: fonts.Cyrillic} {
		_, err := os.Stat("pkg/card_db/images/" + file)
		if (err == nil) != (url != "") {
			t.Fatalf("%s: %q does not match the file (%v)", file, url, err)
		}
		if url == "" {
			continue
		}
		rec := httptest.NewRecorder()
		h := cardImages(files)(http.NotFoundHandler())
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+url, nil))
		if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "font/ttf" || !bytes.HasPrefix(rec.Body.Bytes(), []byte{0, 1, 0, 0}) {
			t.Fatalf("%s: %d %q, not a TrueType font", url, rec.Code, rec.Header().Get("Content-Type"))
		}
	}
}
