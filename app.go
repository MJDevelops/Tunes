package main

import (
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
)

// App struct
type App struct {
	PlayingQueue *audio.PlayingQueue

	// Wails context
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.PlayingQueue = &audio.PlayingQueue{}

	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
