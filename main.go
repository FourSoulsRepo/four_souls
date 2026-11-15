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

	"github.com/FourSoulsRepo/four_souls/bindings/about"
	"github.com/FourSoulsRepo/four_souls/bindings/assets"
	"github.com/FourSoulsRepo/four_souls/bindings/host"
	"github.com/FourSoulsRepo/four_souls/bindings/join"
	"github.com/FourSoulsRepo/four_souls/bindings/prefs"
	"github.com/FourSoulsRepo/four_souls/bindings/records"
	"github.com/FourSoulsRepo/four_souls/bindings/setup"
	"github.com/FourSoulsRepo/four_souls/bindings/updater"
	appfiles "github.com/FourSoulsRepo/four_souls/internal/assets"
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
	app := NewApp()
	cardFiles := appfiles.CardsFS()
	hosting := host.New(appPath("records"), appPath("saves"))

	err = wails.Run(&options.App{
		Title:  title,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     frontend,
			Middleware: assets.Middleware(cardFiles),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		// The saved window mode applies from the first frame (ST-01).
		Fullscreen: settings.Load(settingsFile).Fullscreen,
		// Only the bindings/* structs; App is not bound (ADR 030, item 4).
		Bind: []any{
			about.New(), prefs.New(settingsFile), assets.New(cardFiles), setup.New(),
			join.New(app), hosting, records.New(app), updater.New(app),
		},
		OnStartup: app.startup,
		// Ending the context stops a running update (ADR 030, item 5).
		OnShutdown: func(context.Context) {
			app.shutdown()
			_ = hosting.Stop() //nolint:errcheck,contextcheck // closing; Stop has its own limit
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
