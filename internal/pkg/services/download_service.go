package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
	"github.com/wailsapp/wails/v3/pkg/application"
	wailsevents "github.com/wailsapp/wails/v3/pkg/events"
)

type DownloadService struct {
	ctx     context.Context
	queue   *ytdlp.Queue
	queries *db.Queries
	ytDlp   *ytdlp.YtDlp
	app     *application.App
}

type DownloadServiceOptions struct {
	Workers   uint
	Queries   *db.Queries
	Downloads []ytdlp.Download
	Window    *application.WebviewWindow
}

func NewDownloadService(options DownloadServiceOptions) *DownloadService {
	service := &DownloadService{
		ctx:     nil,
		queue:   ytdlp.NewQueue(options.Workers, options.Downloads...),
		queries: options.Queries,
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

func (s *DownloadService) PendingDownloads() []ytdlp.Download {
	var qDownloads []ytdlp.Download
	ctx := context.Background()
	downloads, _ := s.queries.GetPendingDownloads(ctx)

	for i := range downloads {
		var opts []string
		json.Unmarshal([]byte(downloads[i].Options), &opts)
		download, err := s.ytDlp.NewDownload(downloads[i].ID, downloads[i].Url, opts...)
		if err != nil {
			log.Println(err)
			continue
		}
		qDownloads = append(qDownloads, download)
	}

	return qDownloads
}

func (s *DownloadService) EnqueueDownload(url string, opts ...string) (id string) {
	down, _ := s.ytDlp.NewDownload("", url, opts...)

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

func (s *DownloadService) saveQueueState(downloads []*ytdlp.Download) {
	ctx := context.Background()

	for _, d := range downloads {
		_, err := s.queries.GetDownload(ctx, d.ID)
		if err != nil {
			options, _ := json.Marshal(d.Options)
			s.queries.InsertDownload(ctx, db.InsertDownloadParams{ID: d.ID, Options: string(options), FinishedAt: sql.NullTime{}})
		}
	}
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

	err := s.queries.UpdateDownloadFinishedAt(ctx, db.UpdateDownloadFinishedAtParams{
		ID:         download.ID,
		FinishedAt: t,
	})

	if err != nil {
		// The download wasn't created before, create it now
		options, _ := json.Marshal(download.Options)
		s.queries.InsertDownload(ctx, db.InsertDownloadParams{ID: download.ID, FinishedAt: t, Options: string(options)})
	}
}
