package main

import (
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// App holds the methods the frontend calls through Wails bindings.
type App struct {
	settings string // the settings file
}

// NewApp creates the bindings object.
func NewApp() *App {
	return &App{settings: appPath("settings.json")}
}

// Versions returns the app and rules engine versions for the main menu.
func (a *App) Versions() version.Info {
	return version.Get()
}

// Notice returns the fan-game notice for the splash screen.
func (a *App) Notice() []legal.Line {
	return legal.Notice
}

// Settings returns the saved settings, or the defaults.
func (a *App) Settings() settings.Settings {
	return settings.Load(a.settings)
}

// SaveSettings stores the settings and returns what was stored.
func (a *App) SaveSettings(s settings.Settings) (settings.Settings, error) {
	return settings.Save(a.settings, s) //nolint:wrapcheck // shown to the player as is
}
