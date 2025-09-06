package main

import (
	"context"
	"embed"
	"log"
	"sync"

	"github.com/google/uuid"
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
	ydl, _ := ytdlp.GetLatestRelease()

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
		log.Fatalf("Couldn't initialize connection to database: %v", err)
	}
	db.Migrate()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Tunes",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.SetContext(ctx)
			ytdlp.Initialize(ydl, ctx, db)
			go ydl.StartQueue(queueContext, &queueWg)
			pq.SetContext(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			cancel()
			queueWg.Wait()
			ydl.StopQueue()
		},
		Bind: []interface{}{
			app,
			ydl,
			pq,
		},
		EnumBind: []interface{}{
			events.Events,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: uuid.NewString(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
