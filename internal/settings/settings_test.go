package settings

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "FourSouls", "settings.json")
	saved, err := Save(path, Settings{Nickname: "  Ann  ", Recent: []string{"10.0.0.2"}})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Nickname != "Ann" {
		t.Fatalf("saved nickname %q, want Ann", saved.Nickname)
	}
	// A restart only has the file.
	if got := Load(path); !reflect.DeepEqual(got, saved) {
		t.Fatalf("loaded %+v, want %+v", got, saved)
	}
}

func TestDefaults(t *testing.T) {
	dir := t.TempDir()
	want := Settings{Recent: []string{}, Animation: "normal"}
	if got := Load(filepath.Join(dir, "none.json")); !reflect.DeepEqual(got, want) {
		t.Fatalf("missing file: %+v", got)
	}
	broken := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(broken, []byte("{nick"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(broken); !reflect.DeepEqual(got, want) {
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

func TestRecentIsCleanedAndCapped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	got, err := Save(path, Settings{Recent: []string{" a ", "", "b", "c", "d", "e", "f"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a", "b", "c", "d", "e"}; !reflect.DeepEqual(got.Recent, want) {
		t.Fatalf("recent %q, want %q", got.Recent, want)
	}
}

func TestAnimationSpeed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	got, err := Save(path, Settings{Animation: "fast"})
	if err != nil || got.Animation != "fast" {
		t.Fatalf("fast: %+v, %v", got, err)
	}
	if _, err := Save(path, Settings{Animation: "warp"}); !errors.Is(err, ErrAnimation) {
		t.Fatalf("warp: err %v, want ErrAnimation", err)
	}
	if got := Load(path).Animation; got != "fast" {
		t.Fatalf("animation %q after a bad save, want fast", got)
	}
	if err := os.WriteFile(path, []byte(`{"animation":"warp"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Animation; got != "normal" {
		t.Fatalf("unknown speed reads as %q, want normal", got)
	}
}
