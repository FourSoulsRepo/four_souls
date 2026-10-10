package main

import (
	"context"
	"io/fs"

	"github.com/FourSoulsRepo/four_souls/internal/assets"
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// App holds the methods the frontend calls through Wails bindings.
type App struct {
	ctx       context.Context //nolint:containedctx // Wails hands the app context over at startup
	settings  string          // the settings file
	cardFiles fs.FS           // card data and images
	cardList  cardList
}

// NewApp creates the bindings object.
func NewApp() *App {
	return &App{settings: appPath("settings.json"), cardFiles: assets.CardsFS()}
}

// startup keeps the app context for runtime calls such as dialogs.
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

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

// CardSets returns the codes of the card sets this app has; the client
// sends them when it connects (CD-03).
func (a *App) CardSets() []string {
	var codes []string
	for _, s := range cards.Sets() {
		codes = append(codes, s.Code)
	}
	return codes
}
