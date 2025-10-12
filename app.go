package main

import (
	"context"
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/ffmpeg"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Application state
type App struct {
	YtDlp           *ytdlp.YtDlp
	Ffmpeg          *ffmpeg.Ffmpeg
	PlayingQueue    *audio.Queue
	YtDownloadQueue *ytdlp.Queue
	db              *db.DB
	config          config.Application
	ctx             context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.PlayingQueue = &audio.Queue{}

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	var wg sync.WaitGroup

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	conn.Migrate()
	a.db = conn

	config, err := config.LoadApplicationConfig(path.Join(".", "tunes.config.json"))
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	a.config = config

	wg.Go(func() {
		ytdlp, err := ytdlp.DownloadLatest(path.Join(".", "bin"))
		if err != nil {
			log.Fatalf("Error fetching latest yt-dlp release: %v\n", err)
		}
		ytdlpAbs, _ := filepath.Abs(ytdlp.Path)
		config.Executables.YtDlp.Path = ytdlpAbs
		config.Executables.YtDlp.Release = ytdlp.Release
		a.YtDlp = ytdlp
	})

	wg.Go(func() {
		a.Ffmpeg = ffmpeg.NewFfmpeg()
		err = a.Ffmpeg.DownloadLatest()
		if err != nil {
			log.Fatalf("Error fetching ffmpeg: %v\n", err)
		}

		ffmpegAbs, _ := filepath.Abs(a.Ffmpeg.Path)

		config.Ffmpeg.Version = a.Ffmpeg.Version()
		config.Ffmpeg.Path = ffmpegAbs
	})

	wg.Wait()
	config.Write()
}

func (a *App) initialize() {
	// Load all pending downloads from database
	downloads := a.PendingDownloads()
	a.YtDownloadQueue = ytdlp.NewQueue(5, downloads...).OnShutdown(a.saveQueueState)
	a.YtDownloadQueue.Start()
}

func (a *App) beforeClose(ctx context.Context) bool {
	if a.YtDownloadQueue.IsRunning() {
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
	a.YtDownloadQueue.Stop()
}

func (a *App) EventsEmit(event events.Event, optionalData ...any) {
	runtime.EventsEmit(a.ctx, string(event), optionalData...)
}
