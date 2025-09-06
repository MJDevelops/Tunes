package main

import (
	"context"
	"embed"
	"log"
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

	// Context and waitgroup for the download queue
	var queueWg sync.WaitGroup
	ctx := context.Background()
	queueContext, cancel := context.WithCancel(ctx)

	// Create an instance of the app structure
	app := NewApp()
	pq := &sound.PlayingQueue{}

	// Initialize db connection
	db, err := db.NewDB()
	if err != nil {
		log.Fatalf("Couldn't initialize connection to databaser: %v", err)
	}
	defer db.Close()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Tunes-Gui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.SetContext(ctx)
			ytdlp.Initialize(ctx, db)
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
