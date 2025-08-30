package main

import (
	"context"
	"embed"

	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp/download"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Fetch latest ytdlp version
	ytdlp, _ := ytdlp.GetLatestRelease()

	// Create an instance of the app structure
	app := NewApp()

	// Create download queue
	dq := &download.DownloadQueue{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Tunes-Gui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.SetContext(ctx)
			ytdlp.SetContext(ctx)
			dq.SetContext(ctx)
		},
		Bind: []interface{}{
			app,
			ytdlp,
			dq,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
