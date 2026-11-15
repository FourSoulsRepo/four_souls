package prefs

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// TestSettingsSurviveRestart: a new binding reads what the old one saved.
func TestSettingsSurviveRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := New(path).SaveSettings(settings.Settings{Nickname: "Ann"}); err != nil {
		t.Fatal(err)
	}
	if got := New(path).Settings().Nickname; got != "Ann" {
		t.Fatalf("nickname %q after restart, want Ann", got)
	}
}

// A refused setting comes back in the player's words (ADR 021).
func TestSaveSettingsMessages(t *testing.T) {
	p := New(filepath.Join(t.TempDir(), "settings.json"))
	for _, c := range []struct {
		name  string
		s     settings.Settings
		cause error
		key   string
	}{
		{"nickname", settings.Settings{Nickname: " "}, settings.ErrNickname, "error.nickname_length"},
		{"volume", settings.Settings{Nickname: "Ann", Volume: -1}, settings.ErrChoice, "error.setting_value"},
	} {
		_, err := p.SaveSettings(c.s)
		var m msgerr.Error
		if !errors.As(err, &m) || m.Msg.Key.String() != c.key || !errors.Is(err, c.cause) {
			t.Errorf("%s: %v, want %s", c.name, err, c.key)
		}
	}
}
