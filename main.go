// Command four_souls is the Four Souls desktop app.
package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

//go:embed all:frontend/dist
var frontend embed.FS

func main() {
	app := NewApp()
	host := &Host{}

	err := wails.Run(&options.App{
		Title:  "Four Souls",
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
		OnShutdown: func(ctx context.Context) { _ = host.stop(ctx) }, //nolint:errcheck // the app is closing
	})
	if err != nil {
		log.Fatal(err)
	}
}
