package main

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/FourSoulsRepo/four_souls/internal/assets"
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
	// Russian names come from the old foursouls.ru site (16.5).
	if isaac.Names["en"] != "Isaac" || isaac.Names["ru"] != "Исаак" || len(isaac.Translators["ru"]) == 0 {
		t.Fatalf("isaac names %v, translators %v", isaac.Names, isaac.Translators)
	}
	// The private images submodule may be missing; then no card has an
	// image URL and every card shows as a blank card.
	if _, err := os.Stat("pkg/card_db/images/en/b2/isaac.webp"); (err == nil) != (isaac.Image == "/cardimg/en/b2/isaac.webp") {
		t.Fatalf("isaac image %q does not match the file (%v)", isaac.Image, err)
	}
}

// A two-sided card's faces are two cards, each with its own name, linked
// both ways (ADR 023).
func TestCardsTwoSided(t *testing.T) {
	a := &App{cardFiles: os.DirFS("pkg/card_db")}
	cards, err := a.Cards()
	if err != nil {
		t.Fatal(err)
	}
	byRef := map[string]CardInfo{}
	for _, c := range cards {
		byRef[c.Ref] = c
	}
	front, back := byRef["anima_sola"], byRef["the_revenant"]
	if front.Name != "Anima Sola" || front.Back != "the_revenant" || front.Front != "" {
		t.Errorf("front: %+v", front)
	}
	if back.Name != "The Revenant" || back.Front != "anima_sola" || back.Names["ru"] != "Мститель" || len(back.Stats) != 3 {
		t.Errorf("back: %+v", back)
	}
	if _, err := os.Stat("pkg/card_db/images/en/r/the_revenant.webp"); (err == nil) != (back.Image == "/cardimg/en/r/the_revenant.webp") {
		t.Errorf("back image %q does not match the file (%v)", back.Image, err)
	}
}

func TestCardImages(t *testing.T) {
	files := fstest.MapFS{
		"images/en/b2/isaac.webp":         {Data: []byte("RIFF")},
		"images/assets/die.webp":          {Data: []byte("RIFF")},
		"images/assets/die.png":           {Data: []byte("PNG")},
		"images/assets/fonts/title.ttf":   {Data: []byte("\x00\x01\x00\x00")},
		"images/assets/fonts/LICENSE.txt": {Data: []byte("MIT")},
		"data/b2.json":                    {Data: []byte("{}")},
	}
	frontend := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := cardImages(files)(frontend)
	for path, want := range map[string]int{
		"/cardimg/en/b2/isaac.webp":         http.StatusOK,
		"/cardimg/b2/isaac.webp":            http.StatusNotFound,
		"/cardimg/en/b2/cain.webp":          http.StatusNotFound,
		"/cardimg/../data/b2.json":          http.StatusNotFound,
		"/cardimg/b2/":                      http.StatusNotFound,
		"/index.html":                       http.StatusTeapot,
		"/cardimgs/en/b2/isaac.webp":        http.StatusTeapot,
		"/cardimg/en/b2/../b2/isaac.webp":   http.StatusNotFound,
		"/cardimg/assets/die.webp":          http.StatusOK,
		"/cardimg/assets/die.png":           http.StatusNotFound,
		"/cardimg/assets/die.jpg":           http.StatusNotFound,
		"/cardimg/assets/fonts/title.ttf":   http.StatusOK,
		"/cardimg/assets/fonts/LICENSE.txt": http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+path, nil))
		if rec.Code != want {
			t.Errorf("%s: %d, want %d", path, rec.Code, want)
		}
	}
}

// Each language with an image folder gets its own URL, or the English
// one when the card has no image there (16.6).
func TestCardsImagesPerLanguage(t *testing.T) {
	set := `{"set":{"code":"b2","name":"B","order":1},"cards":[
		{"id":"isaac","type":"Character Card","copies":1,"text":{"en":{"name":"Isaac"}},"images":["b2/isaac.webp"]},
		{"id":"cain","type":"Character Card","copies":1,"text":{"en":{"name":"Cain"}},"images":["b2/cain.webp"]},
		{"id":"maggy","type":"Character Card","copies":1,"text":{"en":{"name":"Maggy"}},"images":["b2/maggy.webp"]}]}`
	a := &App{cardFiles: fstest.MapFS{
		"data/b2.json":            {Data: []byte(set)},
		"images/en/b2/isaac.webp": {Data: []byte("RIFF")},
		"images/en/b2/cain.webp":  {Data: []byte("RIFF")},
		"images/ru/b2/isaac.webp": {Data: []byte("RIFF")},
		"images/assets/die.webp":  {Data: []byte("RIFF")},
	}}
	cards, err := a.Cards()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]string{
		"isaac": {"en": "/cardimg/en/b2/isaac.webp", "ru": "/cardimg/ru/b2/isaac.webp"},
		"cain":  {"en": "/cardimg/en/b2/cain.webp", "ru": "/cardimg/en/b2/cain.webp"},
		"maggy": {},
	}
	for _, c := range cards {
		if !maps.Equal(c.Images, want[c.Ref]) || c.Image != want[c.Ref]["en"] {
			t.Errorf("%s: image %q, images %v; want %v", c.Ref, c.Image, c.Images, want[c.Ref])
		}
	}
}

// A language with card names only is sent in names, so the client lists
// it in Settings without an interface file (16.10).
func TestCardsNamesCardOnlyLanguage(t *testing.T) {
	set := `{"set":{"code":"b2","name":"B","order":1},"cards":[
		{"id":"isaac","type":"Character Card","copies":1,"text":{"en":{"name":"Isaac"},"xx":{"name":"Xisaac"}}},
		{"id":"cain","type":"Character Card","copies":1,"text":{"en":{"name":"Cain"}}}]}`
	a := &App{cardFiles: fstest.MapFS{"data/b2.json": {Data: []byte(set)}}}
	cards, err := a.Cards()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]string{
		"isaac": {"en": "Isaac", "xx": "Xisaac"},
		"cain":  {"en": "Cain"},
	}
	for _, c := range cards {
		if !maps.Equal(c.Names, want[c.Ref]) {
			t.Errorf("%s: names %v, want %v", c.Ref, c.Names, want[c.Ref])
		}
	}
}

// A translator's own image folder (16.12) reaches the client through the
// real card files: App.Cards names it and /cardimg/ serves it, with or
// without the private images.
func TestCardsOwnImages(t *testing.T) {
	own := t.TempDir()
	if err := os.MkdirAll(filepath.Join(own, "de", "b2"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(own, "de", "b2", "isaac.webp"), []byte("RIFF own"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(assets.EnvCardsDir, "pkg/card_db")
	t.Setenv(assets.EnvImagesDir, own)
	a := &App{cardFiles: assets.CardsFS()}
	cards, err := a.Cards()
	if err != nil {
		t.Fatal(err)
	}
	byRef := map[string]CardInfo{}
	for _, c := range cards {
		byRef[c.Ref] = c
	}
	if got := byRef["isaac"].Images["de"]; got != "/cardimg/de/b2/isaac.webp" {
		t.Errorf("isaac in de: %q", got)
	}
	// Cain has no own image: the English one, or none without the private images.
	if got, en := byRef["cain"].Images["de"], byRef["cain"].Image; got != en {
		t.Errorf("cain in de: %q, want the English %q", got, en)
	}
	rec := httptest.NewRecorder()
	cardImages(a.cardFiles)(http.NotFoundHandler()).ServeHTTP(rec,
		httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost/cardimg/de/b2/isaac.webp", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "RIFF own" || rec.Header().Get("Content-Type") != "image/webp" {
		t.Errorf("GET own image: %d %q %q", rec.Code, rec.Body.String(), rec.Header().Get("Content-Type"))
	}
}
