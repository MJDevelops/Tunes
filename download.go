package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
	"gorm.io/gorm"
)

func (a *App) saveQueueState(downloads []ytdlp.Download) {
	ctx := context.Background()
	var dbDownloads []db.Download
	g := gorm.G[db.Download](a.db.Conn())

	for _, d := range downloads {
		_, err := g.Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	g.CreateInBatches(ctx, &dbDownloads, 10)
}

func (a *App) finishDownload(download *ytdlp.Download) {
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

func (a *App) PendingDownloads() []ytdlp.Download {
	var qDownloads []ytdlp.Download
	ctx := context.Background()
	conn := a.db.Conn()
	downloads, _ := gorm.G[db.Download](conn).Where("finished_at IS NULL").Find(ctx)
	for _, d := range downloads {
		qDownloads = append(qDownloads, ytdlp.Download{ID: d.ID, Url: d.Url})
	}

	return qDownloads
}

func (a *App) EnqueueDownload(url string, opts ...string) (id string) {
	down := a.YtDlp.NewDownload(url, opts...)

	down.OnFinished(func() {
		a.finishDownload(&down)
		a.EventsEmit(events.DownloadFinished, down.ID)
	})

	down.OnProgress(func(pf ytdlp.ProgressFormat) {
		a.EventsEmit(events.DownloadProgress, down.ID, pf)
	})

	down.OnStart(func() {
		a.EventsEmit(events.DownloadStarted, down.ID)
	})

	a.YtDownloadQueue.Enqueue(down)

	return down.ID
}
