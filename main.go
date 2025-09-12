package main

import (
	"context"
	"embed"
	"log"
	"sync"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	var err error
	var wg sync.WaitGroup

	aCtx := context.Background()
	appCtx, cancel := context.WithCancel(aCtx)

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	conn.Migrate()

	// Initialize yt-dlp
	ydl, err := ytdlp.Initialize(appCtx, &wg, conn)
	if err != nil {
		log.Fatalf("Error initializing yt-dlp: %v\n", err)
	}

	// Create an instance of the app structure
	app := NewApp()

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
			app.startup(ctx)
			ydl.SetContext(ctx)
			go ydl.StartQueue()
		},
		OnShutdown: func(_ context.Context) {
			cancel()
			wg.Wait()
			ydl.StopQueue()
		},
		Bind: []interface{}{
			app,
			ydl,
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
