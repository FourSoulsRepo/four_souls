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
		"images/" + attackHitSound:    {Data: []byte("ID3")},
		"images/" + effectCastSound:   {Data: []byte("ID3")},
		"images/" + effectImpactSound: {Data: []byte("ID3")},
	}}
	got := a.Sounds()
	if got.AttackHit != "/cardimg/"+attackHitSound || got.EffectCast != "/cardimg/"+effectCastSound || got.EffectImpact != "/cardimg/"+effectImpactSound {
		t.Fatalf("%+v", got)
	}
	// Without the private images the game is silent.
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
	for file, url := range map[string]string{attackHitSound: sounds.AttackHit, effectCastSound: sounds.EffectCast, effectImpactSound: sounds.EffectImpact} {
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
