package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/fstest"
)

func TestCards(t *testing.T) {
	a := &App{cardFiles: os.DirFS("pkg/card_db")}
	cards, err := a.Cards()
	if err != nil {
		t.Fatal(err)
	}
	byRef := map[string]CardInfo{}
	for _, c := range cards {
		byRef[c.Ref] = c
	}
	isaac, ok := byRef["isaac"]
	if !ok || isaac.Name != "Isaac" || isaac.Type != "Character Card" || len(isaac.Stats) != 2 {
		t.Fatalf("isaac: %+v", isaac)
	}
	// The private images submodule may be missing; then no card has an
	// image URL and every card shows as a blank card.
	if _, err := os.Stat("pkg/card_db/images/b2/isaac.webp"); (err == nil) != (isaac.Image == "/cardimg/b2/isaac.webp") {
		t.Fatalf("isaac image %q does not match the file (%v)", isaac.Image, err)
	}
}

func TestCardImages(t *testing.T) {
	files := fstest.MapFS{
		"images/b2/isaac.webp": {Data: []byte("RIFF")},
		"data/b2.json":         {Data: []byte("{}")},
	}
	frontend := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := cardImages(files)(frontend)
	for path, want := range map[string]int{
		"/cardimg/b2/isaac.webp":       http.StatusOK,
		"/cardimg/b2/cain.webp":        http.StatusNotFound,
		"/cardimg/../data/b2.json":     http.StatusNotFound,
		"/cardimg/b2/":                 http.StatusNotFound,
		"/index.html":                  http.StatusTeapot,
		"/cardimgs/b2/isaac.webp":      http.StatusTeapot,
		"/cardimg/b2/../b2/isaac.webp": http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+path, nil))
		if rec.Code != want {
			t.Errorf("%s: %d, want %d", path, rec.Code, want)
		}
	}
}
