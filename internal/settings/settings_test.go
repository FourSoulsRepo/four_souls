package settings

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "FourSouls", "settings.json")
	saved, err := Save(path, Settings{Nickname: "  Ann  "})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Nickname != "Ann" {
		t.Fatalf("saved nickname %q, want Ann", saved.Nickname)
	}
	// A restart only has the file.
	if got := Load(path); got != saved {
		t.Fatalf("loaded %+v, want %+v", got, saved)
	}
}

func TestDefaults(t *testing.T) {
	dir := t.TempDir()
	if got := Load(filepath.Join(dir, "none.json")); got != (Settings{}) {
		t.Fatalf("missing file: %+v", got)
	}
	broken := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(broken, []byte("{nick"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(broken); got != (Settings{}) {
		t.Fatalf("broken file: %+v", got)
	}
}

func TestBadNicknameKeepsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := Save(path, Settings{Nickname: "Ann"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"   ", strings.Repeat("x", 21), "a\nb"} {
		if _, err := Save(path, Settings{Nickname: name}); !errors.Is(err, ErrNickname) {
			t.Fatalf("%q: err %v, want ErrNickname", name, err)
		}
	}
	if got := Load(path).Nickname; got != "Ann" {
		t.Fatalf("nickname %q after bad saves, want Ann", got)
	}
}
