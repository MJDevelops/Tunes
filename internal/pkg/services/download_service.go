package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db/models"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
	"github.com/wailsapp/wails/v3/pkg/application"
	wailsevents "github.com/wailsapp/wails/v3/pkg/events"
	"gorm.io/gorm"
)

type DownloadService struct {
	ctx   context.Context
	queue *ytdlp.Queue
	db    *gorm.DB
	ytDlp *ytdlp.YtDlp
	app   *application.App
}

type DownloadServiceOptions struct {
	Workers   uint
	Db        *gorm.DB
	Downloads []ytdlp.Download
	Window    *application.WebviewWindow
}

func NewDownloadService(options DownloadServiceOptions) *DownloadService {
	service := &DownloadService{
		ctx:   nil,
		queue: ytdlp.NewQueue(options.Workers, options.Downloads...),
		db:    options.Db,
	}

	if options.Window != nil {
		options.Window.RegisterHook(wailsevents.Common.WindowClosing, service.closeHook)
	}

	service.app = application.Get()

	return service
}

func (s *DownloadService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.ctx = ctx
	s.queue.OnShutdown(s.saveQueueState)
	s.queue.Start()
	return nil
}

func (s *DownloadService) ServiceShutdown() error {
	s.queue.Stop()
	return nil
}

func (s *DownloadService) PendingDownloads() ([]ytdlp.Download, error) {
	var qDownloads []ytdlp.Download
	ctx := context.Background()
	downloads, err := gorm.G[models.Download](s.db).Where("finished_at IS NULL").Find(ctx)
	if err != nil {
		return nil, err
	}

	for i := range downloads {
		var opts []string
		json.Unmarshal([]byte(downloads[i].Options), &opts)
		download, err := s.ytDlp.NewDownload(&ytdlp.DownloadOptions{
			ID:      downloads[i].ID,
			URL:     downloads[i].Source,
			Options: opts,
		})
		if err != nil {
			log.Println(err)
			continue
		}
		qDownloads = append(qDownloads, download)
	}

	return qDownloads, nil
}

func (s *DownloadService) EnqueueDownload(url string, opts ...string) (id string) {
	down, _ := s.ytDlp.NewDownload(&ytdlp.DownloadOptions{
		URL:     url,
		Options: opts,
	})

	down.OnFinished(func() {
		s.finishDownload(&down)
		s.app.Event.Emit("tunes:dl:finished", down.ID)
	})

	down.OnProgress(func(pf ytdlp.ProgressFormat) {
		s.app.Event.Emit("tunes:dl:progress", events.DownloadProgress{
			ID:   down.ID,
			Data: pf,
		})
	})

	down.OnStart(func() {
		s.app.Event.Emit("tunes:dl:started", down.ID)
	})

	s.queue.Enqueue(&down)

	return down.ID
}

func (s *DownloadService) saveQueueState(downloads []*ytdlp.Download) error {
	ctx := context.Background()

	for _, d := range downloads {
		_, err := gorm.G[models.Download](s.db).Where("id = ?", d.ID).First(ctx)
		if err != nil {
			options, _ := json.Marshal(d.Options)
			download := &models.Download{
				ID:         d.ID,
				Options:    string(options),
				FinishedAt: sql.NullTime{Valid: false},
			}
			err := gorm.G[models.Download](s.db).Create(ctx, download)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *DownloadService) closeHook(event *application.WindowEvent) {
	if s.queue.IsRunning() {
		qDialog := s.app.Dialog.Question()
		qDialog.SetTitle("Downloads are running")
		qDialog.SetMessage("Downloads are still running. Are you sure you want to quit?")
		qDialog.AddButton("Yes")
		noBtn := qDialog.AddButton("No")
		noBtn.OnClick(func() {
			event.Cancel()
		})
		qDialog.SetDefaultButton(noBtn)
		qDialog.Show()
	}
}

func (s *DownloadService) finishDownload(download *ytdlp.Download) {
	// Try to update existing
	ctx := context.Background()
	t := sql.NullTime{Time: time.Now(), Valid: true}

	_, err := gorm.G[models.Download](s.db).Where("id = ?", download.ID).Update(ctx, "finished_at", t)

	if err != nil {
		// The download wasn't created before, create it now
		options, _ := json.Marshal(download.Options)
		gorm.G[models.Download](s.db).Create(ctx, &models.Download{ID: download.ID, FinishedAt: t, Options: string(options)})
	}
}
