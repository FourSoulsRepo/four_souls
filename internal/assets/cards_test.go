package assets

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
)

func TestCardsFSLoads(t *testing.T) {
	if !Embedded {
		// Tests run in the package folder; point external mode at the module.
		t.Setenv(EnvCardsDir, filepath.Join("..", "..", "pkg", "card_db"))
	}
	if _, err := carddb.Load(CardsFS()); err != nil {
		t.Fatalf("load cards (embedded=%v): %v", Embedded, err)
	}
}

// Images are optional: a fresh clone has none and must still load every card.
func TestCardsFSImagesAreOptional(t *testing.T) {
	if !Embedded {
		t.Setenv(EnvCardsDir, filepath.Join("..", "..", "pkg", "card_db"))
	}
	fsys := CardsFS()
	db, err := carddb.Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	with := 0
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		if carddb.HasImage(fsys, c) {
			with++
		}
	}
	t.Logf("embedded=%v: %d cards, %d with images", Embedded, len(db.Refs()), with)
}

func TestCardsFSEnvOverride(t *testing.T) {
	if Embedded {
		t.Skip("external mode only")
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o750); err != nil {
		t.Fatal(err)
	}
	set := `{"set":{"code":"zz","name":"Test","order":1},"cards":[{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`
	if err := os.WriteFile(filepath.Join(dir, "data", "zz.json"), []byte(set), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvCardsDir, dir)
	db, err := carddb.Load(CardsFS())
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Refs()) != 1 {
		t.Errorf("refs = %v", db.Refs())
	}
}

// FOUR_SOULS_IMAGES lays a translator's own images over the card files
// in every build mode (16.12). Reads never leave that folder.
func TestCardsFSOwnImages(t *testing.T) {
	if !Embedded {
		t.Setenv(EnvCardsDir, filepath.Join("..", "..", "pkg", "card_db"))
	}
	outside := filepath.Join(t.TempDir(), "secret.webp")
	own := t.TempDir()
	write := func(name, data string) {
		t.Helper()
		p := filepath.Join(own, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("de/b2/the_d6.webp", "own de")
	if err := os.WriteFile(outside, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(own, "de", "b2", "isaac.webp")); err != nil {
		t.Skip("no symbolic links here:", err)
	}
	t.Setenv(EnvImagesDir, own)
	fsys := CardsFS()
	if data, err := fs.ReadFile(fsys, "images/de/b2/the_d6.webp"); err != nil || string(data) != "own de" {
		t.Errorf("own image: %q, %v", data, err)
	}
	if data, err := fs.ReadFile(fsys, "images/de/b2/isaac.webp"); err == nil {
		t.Errorf("a link out of the folder was followed: %q", data)
	}
	if _, err := carddb.Load(fsys); err != nil {
		t.Errorf("card data with own images: %v", err)
	}
	if !slices.Contains(carddb.ImageLangs(fsys), "de") {
		t.Errorf("ImageLangs = %v, want de in it", carddb.ImageLangs(fsys))
	}
	// A folder that is not there changes nothing.
	t.Setenv(EnvImagesDir, filepath.Join(own, "missing"))
	if _, err := fs.ReadFile(CardsFS(), "images/de/b2/the_d6.webp"); err == nil {
		t.Error("an image from a missing folder")
	}
}
