package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/fstest"
)

func TestSoundsFromFiles(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{
		"images/" + fireballCastSound:   {Data: []byte("ID3")},
		"images/" + fireballImpactSound: {Data: []byte("ID3")},
	}}
	got := a.Sounds()
	if got.FireballCast != "/cardimg/"+fireballCastSound || got.FireballImpact != "/cardimg/"+fireballImpactSound {
		t.Fatalf("%+v", got)
	}
	// Without the private images the card hit sound is played.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	if got := a.Sounds(); got != (Sounds{}) {
		t.Fatalf("%+v without images", got)
	}
}

// The sounds are MP3 in the private images submodule. With the submodule
// they are found and served as audio; without it, nothing is.
func TestSoundsServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	sounds := (&App{cardFiles: files}).Sounds()
	for file, url := range map[string]string{fireballCastSound: sounds.FireballCast, fireballImpactSound: sounds.FireballImpact} {
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
		if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "audio/mpeg" || rec.Body.Len() == 0 {
			t.Fatalf("%s: %d %q, not an MP3", url, rec.Code, rec.Header().Get("Content-Type"))
		}
	}
}
