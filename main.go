package main

import (
	"embed"
	"log"
	"path"

	"github.com/mjdevelops/tunes/internal/pkg/config"
	tunesdb "github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/services"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var (
		err     error
		binPath = path.Join(".", "bin")
	)

	db, err := gorm.Open(sqlite.Open("tunes.db"), &gorm.Config{})
	tunesdb.Migrate(db)

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
			application.NewService(services.NewYtDlpService(binPath, config)),
			application.NewService(services.NewAudioService(db)),
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
		Db:      db,
		Workers: 5,
		Window:  mainWindow,
	})))

	events.RegisterWailsEvents()

	err = app.Run()

	if err != nil {
		println("Error:", err.Error())
	}
}
