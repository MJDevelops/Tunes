package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/download"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

func (a *App) saveQueueState(downloads <-chan download.Download) {
	ctx := context.Background()
	conn := a.db.Conn()
	var dbDownloads []db.Download

	for d := range downloads {
		_, err := gorm.G[db.Download](conn).Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	gorm.G[db.Download](conn).CreateInBatches(ctx, &dbDownloads, 10)
}

func (a *App) finishDownload(download *download.Download) {
	// Try to update existing
	ctx := context.Background()
	conn := a.db.Conn()
	currTime := time.Now()
	_, err := gorm.G[db.Download](conn).Where("id = ?", download.ID).Update(ctx, "finished_at", currTime)
	if err != nil {
		// The download wasn't created before, create it now
		t := sql.NullTime{Valid: true, Time: currTime}
		dn := db.Download{ID: download.ID, Url: download.Url, FinishedAt: t}
		gorm.G[db.Download](conn).Create(ctx, &dn)
	}
}

func (a *App) loadPendingFromDB() []download.Download {
	var qDownloads []download.Download
	ctx := context.Background()
	conn := a.db.Conn()
	downloads, _ := gorm.G[db.Download](conn).Where("finished_at IS NULL").Find(ctx)
	for _, d := range downloads {
		qDownloads = append(qDownloads, download.Download{ID: d.ID, Url: d.Url})
	}

	return qDownloads
}

func (a *App) EnqueueDownload(url string, opts ...string) (id string) {
	down := download.NewDownload(a.YtDlp.Path, url, opts...)

	down.OnFinished(func() {
		a.finishDownload(&down)
		runtime.EventsEmit(a.ctx, string(events.DownloadFinished), down.ID)
	})

	down.OnProgress(func(pf download.ProgressFormat) {
		runtime.EventsEmit(a.ctx, string(events.DownloadProgress), pf)
	})

	return a.DownloadQueue.SendToQueue(down)
}
