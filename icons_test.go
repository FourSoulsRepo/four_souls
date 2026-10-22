package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
	"testing/fstest"
)

func TestIconsFromFiles(t *testing.T) {
	files := fstest.MapFS{
		"images/assets/icons/coin.webp":  {Data: []byte("RIFF")},
		"images/assets/icons/sword.webp": {Data: []byte("RIFF")},
		"images/assets/banner.webp":      {Data: []byte("RIFF")},
	}
	for _, file := range skyLayers {
		files["images/"+file] = &fstest.MapFile{Data: []byte("RIFF")}
	}
	a := &App{cardFiles: files}
	got := a.Icons()
	if got.Coin != "/cardimg/assets/icons/coin.webp" || got.Sword != "/cardimg/assets/icons/sword.webp" ||
		got.Banner != "/cardimg/assets/banner.webp" {
		t.Fatalf("%+v", got)
	}
	if len(got.Sky) != len(skyLayers) || got.Sky[0] != "/cardimg/assets/background/sky.webp" ||
		got.Sky[5] != "/cardimg/assets/background/stars-5.webp" {
		t.Fatalf("sky %v", got.Sky)
	}
	// A missing star layer leaves the menus plain, not half a sky.
	delete(files, "images/assets/background/stars-3.webp")
	if sky := a.Icons().Sky; len(sky) != 0 {
		t.Fatalf("sky %v with a layer missing", sky)
	}
	// Without the private images the coin and the sword are drawn, the
	// home screen shows its text title, and the menus a plain color.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	if got := a.Icons(); got.Coin != "" || got.Sword != "" || got.Banner != "" || !slices.Equal(got.Sky, []string{}) {
		t.Fatalf("%+v without images", got)
	}
}

// The real coin is WebP in the private images submodule. With the
// submodule it is found and served; without it, the coin is drawn.
func TestCoinServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	coin := (&App{cardFiles: files}).Icons().Coin
	_, err := os.Stat("pkg/card_db/images/" + coinIcon)
	if (err == nil) != (coin != "") {
		t.Fatalf("coin %q does not match the file (%v)", coin, err)
	}
	if coin == "" {
		return
	}
	rec := httptest.NewRecorder()
	h := cardImages(files)(http.NotFoundHandler())
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+coin, nil))
	if rec.Code != http.StatusOK || !bytes.HasPrefix(rec.Body.Bytes(), []byte("RIFF")) {
		t.Fatalf("%s: %d, not a WebP file", coin, rec.Code)
	}
}
