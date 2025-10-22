package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ffmpeg"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

// App Application state
type App struct {
	ytDlp         *ytdlp.YtDlp
	ffmpeg        *ffmpeg.Ffmpeg
	playingQueue  *audio.Queue
	downloadQueue *ytdlp.Queue
	db            *sql.DB
	queries       *db.Queries
	config        config.Application
	ctx           context.Context
}

//go:embed schema.sql
var ddl string

var binPath = path.Join(".", "bin")

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.playingQueue = &audio.Queue{}

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	var wg sync.WaitGroup

	// Initialize db connection
	conn, err := sql.Open("sqlite", "file:tunes.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	a.db = conn

	bCtx := context.Background()
	if _, err := conn.ExecContext(bCtx, ddl); err != nil {
		log.Fatalf("Error creating db tables: %v\n", err)
	}
	a.queries = db.New(conn)

	config, err := config.LoadApplicationConfig(path.Join(".", "tunes.config.json"))
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	a.config = config

	wg.Go(func() {
		ytdlp, err := ytdlp.GetLatest(binPath)
		if err != nil {
			log.Fatalf("Error fetching latest yt-dlp release: %v\n", err)
		}

		ytdlpAbs, _ := filepath.Abs(ytdlp.Path())
		config.YtDlp.Path = ytdlpAbs
		config.YtDlp.Release = ytdlp.Release
		a.ytDlp = ytdlp
	})

	wg.Go(func() {
		a.ffmpeg, err = ffmpeg.NewFfmpeg(binPath)
		if err != nil {
			log.Fatalf("Error initializing ffmpeg: %v\n", err)
		}

		if err = a.ffmpeg.GetLatest(); err != nil {
			log.Fatalf("Error fetching ffmpeg: %v\n", err)
		}

		ffmpegAbs, _ := filepath.Abs(a.ffmpeg.Path())
		config.Ffmpeg.Version = a.ffmpeg.Version()
		config.Ffmpeg.Path = ffmpegAbs
	})

	wg.Wait()
	config.Write()
}

func (a *App) initialize() {
	// Load all pending downloads from database
	downloads := a.PendingDownloads()
	a.downloadQueue = ytdlp.NewQueue(5, downloads...).OnShutdown(a.saveQueueState)
	a.downloadQueue.Start()
}

func (a *App) beforeClose(ctx context.Context) bool {
	if a.downloadQueue.IsRunning() {
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
	a.downloadQueue.Stop()
}

func (a *App) EventsEmit(event events.Event, optionalData ...any) {
	runtime.EventsEmit(a.ctx, string(event), optionalData...)
}
