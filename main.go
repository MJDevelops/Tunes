package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"path"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/services"
	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed schema.sql
var ddl string

func main() {
	var (
		err     error
		binPath = path.Join(".", "bin")
	)

	// Initialize db connection
	conn, err := sql.Open("sqlite", "file:tunes.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}

	// Create db tables
	bCtx := context.Background()
	if _, err := conn.ExecContext(bCtx, ddl); err != nil {
		log.Fatalf("Error creating db tables: %v\n", err)
	}

	queries := db.New(conn)

	config, err := config.LoadApplicationConfig(path.Join(".", "tunes.config.json"))
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	app := application.New(application.Options{
		Name: "Tunes",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(services.NewFfmpegService(binPath, config)),
			application.NewService(services.NewYtDlpService(binPath, config)),
			application.NewService(services.NewAudioService(queries)),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "01993fca-6c97-746f-b747-6c0c12b27e32",
		},
	})

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Tunes",
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	app.RegisterService(application.NewService(services.NewDownloadService(services.DownloadServiceOptions{
		Queries: queries,
		Workers: 5,
		Window:  mainWindow,
	})))

	events.RegisterWailsEvents()

	err = app.Run()

	if err != nil {
		println("Error:", err.Error())
	}
}
