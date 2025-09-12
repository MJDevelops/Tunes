package main

import (
	"embed"
	"log"

	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	var err error

	// Create an instance of the app structure
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Error during initialization of app: %v\n", err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Tunes",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		EnumBind: []interface{}{
			events.Events,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "01993fca-6c97-746f-b747-6c0c12b27e32",
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
