package main

import (
	"context"
	"log"
	"sync"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
)

// Application state
type App struct {
	YtDlp         *ytdlp.YtDlp
	PlayingQueue  *audio.PlayingQueue
	DownloadQueue DownloadQueue
	db            *db.DB

	// App context
	aCtx   context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	// Wails context
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.PlayingQueue = &audio.PlayingQueue{}

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	conn.Migrate()
	app.db = conn

	app.aCtx, app.cancel = context.WithCancel(context.Background())

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) initialize() {
	// Load all pending downloads from database
	a.loadPendingFromDB()

	// Start download queue
	go a.startQueue()
}

func (a *App) shutdown(_ context.Context) {
	a.cancel()
	a.wg.Wait()
	a.stopQueue()
}
