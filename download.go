package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/download"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"gorm.io/gorm"
)

func (a *App) saveQueueState(downloads <-chan download.Download) {
	ctx := context.Background()
	var dbDownloads []db.Download
	g := gorm.G[db.Download](a.db.Conn())

	for d := range downloads {
		_, err := g.Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	g.CreateInBatches(ctx, &dbDownloads, 10)
}

func (a *App) finishDownload(download *download.Download) {
	// Try to update existing
	ctx := context.Background()
	currTime := time.Now()
	g := gorm.G[db.Download](a.db.Conn())

	_, err := g.Where("id = ?", download.ID).Update(ctx, "finished_at", currTime)
	if err != nil {
		// The download wasn't created before, create it now
		t := sql.NullTime{Valid: true, Time: currTime}
		dn := db.Download{ID: download.ID, Url: download.Url, FinishedAt: t}
		g.Create(ctx, &dn)
	}
}

func (a *App) PendingDownloads() []download.Download {
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
		a.EventsEmit(events.DownloadFinished, down.ID)
	})

	down.OnProgress(func(pf download.ProgressFormat) {
		a.EventsEmit(events.DownloadProgress, pf)
	})

	return a.DownloadQueue.SendToQueue(down)
}
