package main

import (
	"context"
	"embed"
	"sync"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/sound"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Fetch latest ytdlp version
	ytdlp, _ := ytdlp.GetLatestRelease()

	var queueWg sync.WaitGroup

	// Create an instance of the app structure
	app := NewApp()
	pq := &sound.PlayingQueue{}

	db := db.NewDB()
	defer db.Close()

	ctx := context.Background()
	queueContext, cancel := context.WithCancel(ctx)

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
			db.SetContext(ctx)
			go ytdlp.StartQueue(queueContext, &queueWg)
			pq.SetContext(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			cancel()
			queueWg.Wait()
			ytdlp.StopQueue()
		},
		Bind: []interface{}{
			app,
			ytdlp,
			db,
			pq,
		},
		EnumBind: []interface{}{
			events.Events,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
