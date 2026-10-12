package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

func TestDiceMatchSettings(t *testing.T) {
	var ids []string
	for _, d := range diceSkins {
		ids = append(ids, d.id)
	}
	if !slices.Equal(ids, settings.Dice) {
		t.Fatalf("dice %v, settings %v", ids, settings.Dice)
	}
}

func TestDiceImage(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{"images/" + diceSkins[0].file: {Data: []byte("RIFF")}}}
	if got := a.Dice()[0].Image; got != "/cardimg/"+diceSkins[0].file {
		t.Fatalf("image %q", got)
	}
	// Without the private images every die is drawn plain.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	for _, d := range a.Dice() {
		if d.Image != "" || d.Name == "" {
			t.Fatalf("%+v without images", d)
		}
	}
}

// The real dice textures are WebP in the private images submodule. With
// the submodule, each one is found and served; without it, dice are plain.
func TestDiceTexturesServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	a := &App{cardFiles: files}
	h := cardImages(files)(http.NotFoundHandler())
	for i, d := range a.Dice() {
		_, err := os.Stat("pkg/card_db/images/" + diceSkins[i].file)
		if (err == nil) != (d.Image != "") {
			t.Fatalf("%s: image %q does not match the file (%v)", d.ID, d.Image, err)
		}
		if d.Image == "" {
			continue
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+d.Image, nil))
		if rec.Code != http.StatusOK || !bytes.HasPrefix(rec.Body.Bytes(), []byte("RIFF")) {
			t.Fatalf("%s: %d, not a WebP file", d.Image, rec.Code)
		}
	}
}
