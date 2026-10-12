package main

import (
	"path/filepath"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// TestSettingsSurviveRestart: a new App reads what the old one saved.
func TestSettingsSurviveRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := (&App{settings: path}).SaveSettings(settings.Settings{Nickname: "Ann"}); err != nil {
		t.Fatal(err)
	}
	if got := (&App{settings: path}).Settings().Nickname; got != "Ann" {
		t.Fatalf("nickname %q after restart, want Ann", got)
	}
}

// TestSecondCopySettings: a second copy keeps the first copy's choices
// but has its own nickname and seat.
func TestSecondCopySettings(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "settings.json")
	if _, err := settings.Save(first, settings.Settings{
		Nickname: "Ann", Theme: "light",
		Rejoin: settings.Rejoin{Address: "localhost", Token: "t1", Game: "g1"},
	}); err != nil {
		t.Fatal(err)
	}
	a, err := instance.Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := instance.Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := slotSettings(a, first); got != first {
		t.Errorf("first copy uses %q, want %q", got, first)
	}
	path := slotSettings(b, first)
	if path == first {
		t.Fatal("second copy shares the first copy's settings")
	}
	s := settings.Load(path)
	if s.Nickname != "" || s.Rejoin.Token != "" || s.Theme != "light" {
		t.Errorf("second copy settings %+v; want light theme, no nickname, no seat", s)
	}
	if settings.Load(first).Rejoin.Token != "t1" {
		t.Error("the first copy lost its seat")
	}
}
