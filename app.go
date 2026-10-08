package main

import (
	"context"

	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Versions returns the app and rules engine versions for the main menu.
func (a *App) Versions() version.Info {
	return version.Get()
}

// Notice returns the fan-game notice for the splash screen.
func (a *App) Notice() []legal.Line {
	return legal.Notice
}
