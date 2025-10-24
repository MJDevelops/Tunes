package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
)

func (a *App) saveQueueState(downloads []*ytdlp.Download) {
	ctx := context.Background()

	for _, d := range downloads {
		_, err := a.queries.GetDownload(ctx, d.ID)
		if err != nil {
			options, _ := json.Marshal(d.Options)
			a.queries.InsertDownload(ctx, db.InsertDownloadParams{ID: d.ID, Options: string(options), FinishedAt: sql.NullTime{}})
		}
	}
}

func (a *App) finishDownload(download *ytdlp.Download) {
	// Try to update existing
	ctx := context.Background()
	t := sql.NullTime{Time: time.Now(), Valid: true}

	err := a.queries.UpdateDownloadFinishedAt(ctx, db.UpdateDownloadFinishedAtParams{
		ID:         download.ID,
		FinishedAt: t,
	})

	if err != nil {
		// The download wasn't created before, create it now
		options, _ := json.Marshal(download.Options)
		a.queries.InsertDownload(ctx, db.InsertDownloadParams{ID: download.ID, FinishedAt: t, Options: string(options)})
	}
}

func (a *App) PendingDownloads() []ytdlp.Download {
	var qDownloads []ytdlp.Download
	ctx := context.Background()
	downloads, _ := a.queries.GetPendingDownloads(ctx)

	for i := range downloads {
		var opts []string
		json.Unmarshal([]byte(downloads[i].Options), &opts)
		download, err := a.ytDlp.NewDownload(downloads[i].ID, downloads[i].Url, opts...)
		if err != nil {
			log.Println(err)
			continue
		}
		qDownloads = append(qDownloads, download)
	}

	return qDownloads
}

func (a *App) EnqueueDownload(url string, opts ...string) (id string) {
	down, _ := a.ytDlp.NewDownload("", url, opts...)

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

	a.downloadQueue.Enqueue(&down)

	return down.ID
}
