package main

import (
	"context"
	"log"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/download"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Application state
type App struct {
	YtDlp         *ytdlp.YtDlp
	PlayingQueue  *audio.PlayingQueue
	DownloadQueue *download.DownloadQueue
	db            *db.DB
	ctx           context.Context
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

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) initialize() {
	// Load all pending downloads from database
	downloads := a.loadPendingFromDB()
	a.DownloadQueue = download.NewDownloadQueue(5, downloads...).OnShutdown(a.saveQueueState)
	a.DownloadQueue.Start()
}

func (a *App) beforeClose(ctx context.Context) bool {
	if a.DownloadQueue.IsRunning() {
		dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.QuestionDialog,
			Title:   "Quit",
			Message: "Are you sure you want to quit? Currently running downloads will be cancelled.",
		})

		if err != nil {
			return false
		}

		return dialog != "Yes"
	}

	return false
}

func (a *App) shutdown(_ context.Context) {
	a.DownloadQueue.Stop()
}
