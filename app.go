package main

import (
	"context"
	"log"
	"sync"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
)

// App struct
type App struct {
	YtDlp        *ytdlp.YtDlp
	PlayingQueue *audio.PlayingQueue

	// Wails context
	ctx context.Context

	// App context
	appCtx context.Context

	// Cancellation function for app context
	cancel func()

	wg sync.WaitGroup
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	aCtx := context.Background()
	appCtx, cancel := context.WithCancel(aCtx)
	a.appCtx = appCtx
	a.cancel = cancel

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Couldn't initialize connection to database: %v", err)
	}
	conn.Migrate()

	pq := &audio.PlayingQueue{}
	pq.SetContext(ctx)
	a.PlayingQueue = pq

	ydl, err := ytdlp.Initialize(appCtx, &a.wg, conn)
	if err != nil {
		log.Fatalf("Error during initialization of yt-dlp: %s", err.Error())
	}
	a.YtDlp = ydl

	ydl.SetContext(ctx)
	go ydl.StartQueue()
}

func (a *App) shutdown(ctx context.Context) {
	a.cancel()
	a.wg.Wait()
	a.YtDlp.StopQueue()
}
