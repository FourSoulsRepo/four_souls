package main

import (
	"path/filepath"
	"testing"

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
