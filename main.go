// Command FourSouls is the Four Souls desktop app.
package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	"github.com/FourSoulsRepo/four_souls/internal/tracing"
)

//go:embed all:frontend/dist
var frontend embed.FS

func main() {
	// Traces are written when the app closes (ADR 024).
	stopTrace, err := tracing.Start(context.Background(), tracing.FromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if serr := stopTrace(); serr != nil {
			log.Print(serr)
		}
	}()
	// Each running copy holds a slot with its own settings, so several
	// players can share one computer.
	title, settingsFile := "Four Souls", appPath("settings.json")
	slot, err := instance.Acquire(filepath.Dir(settingsFile))
	if err != nil {
		log.Printf("%v; sharing the first copy's settings", err)
	} else {
		defer slot.Release() //nolint:errcheck // the app is closing
		settingsFile = slotSettings(slot, settingsFile)
		if slot.N > 1 {
			title = fmt.Sprintf("%s (%d)", title, slot.N)
		}
	}
	app := NewApp(settingsFile)
	host := &Host{}

	err = wails.Run(&options.App{
		Title:  title,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     frontend,
			Middleware: cardImages(app.cardFiles),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		// The saved window mode applies from the first frame (ST-01).
		Fullscreen: settings.Load(app.settings).Fullscreen,
		Bind:       []any{app, host},
		OnStartup:  app.startup,
		OnShutdown: func(ctx context.Context) { _ = host.stop(ctx) }, //nolint:errcheck // the app is closing
	})
	if err != nil {
		log.Fatal(err)
	}
}
