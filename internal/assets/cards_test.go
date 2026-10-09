package assets

import (
	"os"
	"path/filepath"
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

func TestCardsFSEnvOverride(t *testing.T) {
	if Embedded {
		t.Skip("external mode only")
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o750); err != nil {
		t.Fatal(err)
	}
	set := `{"set":{"code":"zz","name":"Test"},"cards":[{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`
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
