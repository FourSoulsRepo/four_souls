package main

import (
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// App holds the methods the frontend calls through Wails bindings.
type App struct{}

// NewApp creates the bindings object.
func NewApp() *App {
	return &App{}
}

// Versions returns the app and rules engine versions for the main menu.
func (a *App) Versions() version.Info {
	return version.Get()
}

// Notice returns the fan-game notice for the splash screen.
func (a *App) Notice() []legal.Line {
	return legal.Notice
}
