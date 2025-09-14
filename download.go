package main

import (
	"context"
	"database/sql"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/db"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type Download struct {
	ID       string
	Url      string
	Progress int
}

type DownloadQueue struct {
	Running []Download
	Waiting []Download
	mu      sync.Mutex
	once    sync.Once
}

// Threads to be used for downloads
var MaxThreads int

func init() {
	if MaxThreads = config.GetMaxThreads(); MaxThreads <= 0 {
		MaxThreads = 5
	}
}

func (dq *DownloadQueue) isRunning() bool {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	return len(dq.Running) > 0 || len(dq.Waiting) > 0
}

// Adds download to queue and returns the corresponding ID
func (a *App) AddToQueue(download Download) string {
	dq := &a.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	id := uuid.NewString()
	download.ID = id

	if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
		runtime.EventsEmit(a.ctx, string(events.DownloadQueueStarted))
	}

	dq.Waiting = append(dq.Waiting, download)
	return id
}

func (a *App) startQueue() {
	dq := &a.DownloadQueue
	dq.once.Do(func() {
		a.wg.Add(1)
		defer a.wg.Done()
		for {
			select {
			case <-a.aCtx.Done():
				return
			default:
				dq.mu.Lock()
				for len(dq.Waiting) > 0 && len(dq.Running) < MaxThreads {
					newDown := dq.Waiting[0]

					go func() {
						a.download(newDown)
						if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
							runtime.EventsEmit(a.ctx, string(events.DownloadQueueDone))
						}
					}()

					dq.Running = append(dq.Running, newDown)
					dq.Waiting = slices.Delete(dq.Waiting, 0, 1)
				}
				dq.mu.Unlock()
			}
		}
	})
}

func (a *App) stopQueue() {
	dq := &a.DownloadQueue
	dq.mu.Lock()
	a.saveQueueState()
}

func (a *App) removeFromQueue(id string) {
	dq := &a.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for i, d := range dq.Running {
		if d.ID == id {
			a.finishDownload(&d)
			dq.Running = slices.Delete(dq.Running, i, i+1)
			return
		}
	}
}

func (a *App) saveQueueState() {
	ctx := context.Background()
	var dbDownloads []db.Download

	allUnfinished := append(a.DownloadQueue.Waiting, a.DownloadQueue.Running...)
	for _, d := range allUnfinished {
		_, err := gorm.G[db.Download](a.db.Conn).Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	gorm.G[db.Download](a.db.Conn).CreateInBatches(ctx, &dbDownloads, 10)
}

func (a *App) download(download Download) {
	a.wg.Add(1)
	defer a.wg.Done()
	cmd := a.YtDlp.CreateCommandQuiet(download.Url, "-P", download.ID)
	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	runtime.EventsEmit(a.ctx, string(events.DownloadStarted), download.ID)

	select {
	case <-a.aCtx.Done():
		runtime.EventsEmit(a.ctx, string(events.DownloadInterrupt), download.ID)
		cmd.Cancel()
		os.RemoveAll(download.ID)
		return
	case <-ch:
		runtime.EventsEmit(a.ctx, string(events.DownloadFinished), download.ID)
		a.removeFromQueue(download.ID)
		return
	}
}

func (a *App) finishDownload(download *Download) {
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

func (a *App) loadPendingFromDB() {
	ctx := context.Background()
	downloads, _ := gorm.G[db.Download](a.db.Conn).Where("finished_at IS NULL").Find(ctx)
	for _, d := range downloads {
		a.DownloadQueue.Waiting = append(a.DownloadQueue.Waiting, Download{ID: d.ID, Url: d.Url})
	}
}
