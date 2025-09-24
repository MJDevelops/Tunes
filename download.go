package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/download"
	"gorm.io/gorm"
)

func (a *App) saveQueueState(downloads <-chan download.Download) {
	ctx := context.Background()
	var dbDownloads []db.Download

	for d := range downloads {
		_, err := gorm.G[db.Download](a.db.Conn).Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	gorm.G[db.Download](a.db.Conn).CreateInBatches(ctx, &dbDownloads, 10)
}

func (a *App) finishDownload(download *download.Download) {
	// Try to update existing
	ctx := context.Background()
	currTime := time.Now()
	_, err := gorm.G[db.Download](a.db.Conn).Where("id = ?", download.ID).Update(ctx, "finished_at", currTime)
	if err != nil {
		// The download wasn't created before, create it now
		t := sql.NullTime{Valid: true, Time: currTime}
		dn := db.Download{ID: download.ID, Url: download.Url, FinishedAt: t}
		gorm.G[db.Download](a.db.Conn).Create(ctx, &dn)
	}
}

func (a *App) loadPendingFromDB() []download.Download {
	var qDownloads []download.Download
	ctx := context.Background()
	downloads, _ := gorm.G[db.Download](a.db.Conn).Where("finished_at IS NULL").Find(ctx)
	for _, d := range downloads {
		qDownloads = append(qDownloads, download.Download{ID: d.ID, Url: d.Url})
	}

	return qDownloads
}
