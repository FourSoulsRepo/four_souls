package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/fstest"
)

func TestBacksFromFiles(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{"images/assets/backs/loot.webp": {Data: []byte("RIFF")}}}
	got := a.Backs()
	if len(got) != len(backKinds) {
		t.Fatalf("backs %+v", got)
	}
	for _, b := range got {
		want := ""
		if b.Kind == "loot" {
			want = "/cardimg/assets/backs/loot.webp"
		}
		if b.Image != want {
			t.Fatalf("%s: image %q, want %q", b.Kind, b.Image, want)
		}
	}
	// Without the private images every back is drawn plain.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	for _, b := range a.Backs() {
		if b.Image != "" {
			t.Fatalf("%+v without images", b)
		}
	}
}

// The real backs are WebP in the private images submodule. With the
// submodule, each one is found and served; without it, backs are plain.
func TestBacksServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	a := &App{cardFiles: files}
	h := cardImages(files)(http.NotFoundHandler())
	for _, b := range a.Backs() {
		_, err := os.Stat("pkg/card_db/images/" + backsDir + "/" + b.Kind + ".webp")
		if (err == nil) != (b.Image != "") {
			t.Fatalf("%s: image %q does not match the file (%v)", b.Kind, b.Image, err)
		}
		if b.Image == "" {
			continue
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+b.Image, nil))
		if rec.Code != http.StatusOK || !bytes.HasPrefix(rec.Body.Bytes(), []byte("RIFF")) {
			t.Fatalf("%s: %d, not a WebP file", b.Image, rec.Code)
		}
	}
}
