// Command four_souls is the Four Souls desktop app.
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var frontend embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Four Souls",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontend,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []any{app},
	})
	if err != nil {
		log.Fatal(err)
	}
}
