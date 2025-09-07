package ytdlp

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
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

// Adds download to queue and returns the corresponding ID
func (y *YtDlp) AddToQueue(download *Download) string {
	dq := &y.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	id := uuid.NewString()
	download.ID = id

	if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
		runtime.EventsEmit(y.ctx, string(events.DownloadQueueStarted))
	}

	dq.Waiting = append(dq.Waiting, *download)
	return id
}

func (y *YtDlp) StartQueue(ctx context.Context, wg *sync.WaitGroup) {
	dq := &y.DownloadQueue
	dq.once.Do(func() {
		// Load all pending downloads from database
		y.loadPendingFromDB()
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				dq.mu.Lock()
				for len(dq.Waiting) > 0 && len(dq.Running) < MaxThreads {
					newDown := dq.Waiting[0]

					go func() {
						y.download(ctx, wg, newDown)
						if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
							runtime.EventsEmit(y.ctx, string(events.DownloadQueueDone))
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

func (y *YtDlp) StopQueue() {
	dq := &y.DownloadQueue
	dq.mu.Lock()
	y.saveQueueState()
}

func (y *YtDlp) removeFromQueue(id string) {
	dq := &y.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for i, d := range dq.Running {
		if d.ID == id {
			y.finishDownload(&d)
			dq.Running = slices.Delete(dq.Running, i, i+1)
			return
		}
	}
}

func (y *YtDlp) saveQueueState() {
	ctx := context.Background()
	var dbDownloads []db.Download

	allUnfinished := append(y.DownloadQueue.Waiting, y.DownloadQueue.Running...)
	for _, d := range allUnfinished {
		_, err := gorm.G[db.Download](y.db.Conn).Where("id = ?", d.ID).First(ctx)
		if err != nil {
			dbDownloads = append(dbDownloads, db.Download{ID: d.ID, Url: d.Url})
		}
	}

	gorm.G[db.Download](y.db.Conn).CreateInBatches(ctx, &dbDownloads, 10)
}

func (y *YtDlp) download(ctx context.Context, wg *sync.WaitGroup, download Download) {
	wg.Add(1)
	defer wg.Done()
	cmd := exec.CommandContext(context.Background(), y.Bin, download.Url, "-P", download.ID)
	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	runtime.EventsEmit(y.ctx, string(events.DownloadStarted), download.ID)

	select {
	case <-ctx.Done():
		runtime.EventsEmit(y.ctx, string(events.DownloadInterrupt), download.ID)
		cmd.Cancel()
		os.RemoveAll(download.ID)
		return
	case <-ch:
		runtime.EventsEmit(y.ctx, string(events.DownloadFinished), download.ID)
		y.removeFromQueue(download.ID)
		return
	}
}

func (y *YtDlp) finishDownload(download *Download) {
	// Try to update existing
	ctx := context.Background()
	currTime := time.Now()
	_, err := gorm.G[db.Download](y.db.Conn).Where("id = ?", download.ID).Update(ctx, "finished_at", currTime)
	if err != nil {
		// The download wasn't created before, create it now
		t := sql.NullTime{Valid: true, Time: currTime}
		dn := db.Download{ID: download.ID, Url: download.Url, FinishedAt: t}
		gorm.G[db.Download](y.db.Conn).Create(ctx, &dn)
	}
}

func (y *YtDlp) loadPendingFromDB() {
	ctx := context.Background()
	downloads, _ := gorm.G[db.Download](y.db.Conn).Where("finished_at IS NULL").Find(ctx)
	for _, d := range downloads {
		y.DownloadQueue.Waiting = append(y.DownloadQueue.Waiting, Download{ID: d.ID, Url: d.Url})
	}
}
