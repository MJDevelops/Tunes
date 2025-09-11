package main

import (
	"context"
	"embed"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/mjdevelops/tunes/internal/pkg/audio"
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
	// App waitgroup
	var wg sync.WaitGroup

	// App context
	ctx := context.Background()
	appCtx, cancel := context.WithCancel(ctx)

	// Create an instance of the app structure
	app := NewApp()
	pq := &audio.PlayingQueue{}

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Couldn't initialize connection to database: %v", err)
	}
	conn.Migrate()

	// Initialize yt-dlp
	ydl, _ := ytdlp.Initialize(appCtx, &wg, conn)

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
			ydl.SetContext(ctx)
			pq.SetContext(ctx)
			go ydl.StartQueue()
		},
		OnShutdown: func(ctx context.Context) {
			cancel()
			wg.Wait()
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
