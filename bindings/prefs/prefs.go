// Package prefs is the binding of the player's app settings (ADR 030).
// It is not named settings: Wails would merge its types with those of
// internal/settings in models.ts.
package prefs

import (
	"errors"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Settings errors the player reads in their own language (ADR 021);
// their words are in frontend/src/locales/en.json.
var (
	errNickname     = engine.NewKey("error.nickname_length")
	errSettingValue = engine.NewKey("error.setting_value")
)

// Prefs loads and saves the settings file.
type Prefs struct {
	file string // the settings file
}

// New makes the binding for the settings file.
func New(settingsFile string) *Prefs {
	return &Prefs{file: settingsFile}
}

// Settings returns the saved settings, or the defaults.
func (p *Prefs) Settings() settings.Settings {
	return settings.Load(p.file)
}

// SaveSettings stores the settings and returns what was stored.
func (p *Prefs) SaveSettings(s settings.Settings) (settings.Settings, error) {
	saved, err := settings.Save(p.file, s)
	return saved, settingsError(err)
}

// settingsError gives a settings error the player's words; any other
// error (the disk) stays as it is.
func settingsError(err error) error {
	var m engine.Msg
	switch {
	case err == nil:
		return nil
	case errors.Is(err, settings.ErrNickname):
		m = errNickname.Msg(engine.N("n", protocol.MaxName))
	case errors.Is(err, settings.ErrChoice):
		m = errSettingValue.Msg()
	default:
		return err
	}
	return msgerr.New(m, err)
}
