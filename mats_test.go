package main

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMatsFromFiles(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{
		"images/assets/mats/basement.webp":  {Data: []byte("RIFF")},
		"images/assets/mats/dark_room.webp": {Data: []byte("RIFF")},
		"images/assets/mats/old.png":        {Data: []byte("PNG")},
		"images/assets/mats/Bad Name.webp":  {Data: []byte("RIFF")},
		"images/assets/mats/sub/x.webp":     {Data: []byte("RIFF")},
	}}
	want := []MatInfo{
		{ID: "basement", Name: "Basement", Image: "/cardimg/assets/mats/basement.webp"},
		{ID: "dark_room", Name: "Dark room", Image: "/cardimg/assets/mats/dark_room.webp"},
	}
	if got := a.Mats(); !reflect.DeepEqual(got, want) {
		t.Fatalf("mats %+v, want %+v", got, want)
	}
	// Without the private images the list is empty, never nil.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	if got := a.Mats(); got == nil || len(got) != 0 {
		t.Fatalf("mats %#v without images", got)
	}
}

// The mat list matches the folder in the private images submodule, and
// each mat is served. Without the submodule both are empty.
func TestMatsMatchFolder(t *testing.T) {
	var want []string
	entries, err := os.ReadDir("pkg/card_db/images/assets/mats")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		t.Fatal(err)
	}
	for _, e := range entries {
		if id, ok := strings.CutSuffix(e.Name(), ".webp"); ok {
			want = append(want, id)
		}
	}
	files := os.DirFS("pkg/card_db")
	a := &App{cardFiles: files}
	h := cardImages(files)(http.NotFoundHandler())
	var got []string
	for _, m := range a.Mats() {
		got = append(got, m.ID)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+m.Image, nil))
		if rec.Code != http.StatusOK || !bytes.HasPrefix(rec.Body.Bytes(), []byte("RIFF")) {
			t.Fatalf("%s: %d, not a WebP file", m.Image, rec.Code)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mats %v, folder %v", got, want)
	}
}
