package main

import (
	"context"
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
func NewApp() (*App, error) {
	app := &App{}
	app.PlayingQueue = &audio.PlayingQueue{}

	aCtx := context.Background()
	appCtx, cancel := context.WithCancel(aCtx)
	app.appCtx = appCtx
	app.cancel = cancel

	// Initialize db connection
	conn, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	conn.Migrate()

	ydl, err := ytdlp.Initialize(appCtx, &app.wg, conn)
	if err != nil {
		return nil, err
	}
	app.YtDlp = ydl

	return app, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.PlayingQueue.SetContext(ctx)
	a.YtDlp.SetContext(ctx)
	go a.YtDlp.StartQueue()
}

func (a *App) shutdown(ctx context.Context) {
	a.cancel()
	a.wg.Wait()
	a.YtDlp.StopQueue()
}
