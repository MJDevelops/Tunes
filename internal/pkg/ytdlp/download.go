package ytdlp

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/events"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Download struct {
	ID       string
	Url      string
	Progress int
}

type DownloadQueue struct {
	Running []Download
	Waiting []Download
	rMu     sync.Mutex
	wMu     sync.Mutex
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
	dq.wMu.Lock()
	defer dq.wMu.Unlock()
	id := uuid.NewString()
	download.ID = id

	if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
		runtime.EventsEmit(y.ctx, string(events.DownloadQueueStarted))
	}

	dq.Waiting = append(dq.Waiting, *download)
	return id
}

func (y *YtDlp) StartQueue(ctx context.Context) {
	dq := &y.DownloadQueue
	dq.once.Do(func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				dq.rMu.Lock()
				dq.wMu.Lock()
				if len(dq.Running) < MaxThreads && len(dq.Waiting) > 0 {
					for len(dq.Waiting) > 0 {
						if len(dq.Running) == MaxThreads {
							break
						}

						newDown := dq.Waiting[0]

						go func() {
							y.download(ctx, newDown)
							if len(dq.Running) == 0 && len(dq.Waiting) == 0 {
								runtime.EventsEmit(y.ctx, string(events.DownloadQueueDone))
							}
						}()

						dq.Running = append(dq.Running, newDown)
						dq.Waiting = slices.Delete(dq.Waiting, 0, 0)
					}
				}
				dq.rMu.Unlock()
				dq.wMu.Unlock()
			}
		}
	})
}

func (y *YtDlp) StopQueue() {
	dq := &y.DownloadQueue
	dq.rMu.Lock()
	dq.wMu.Lock()
	y.saveQueueState()
}

func (y *YtDlp) removeFromQueue(id string) {
	dq := &y.DownloadQueue
	dq.rMu.Lock()
	defer dq.rMu.Unlock()
	for i, d := range dq.Running {
		if d.ID == id {
			dq.Running[i] = dq.Running[len(dq.Running)-1]
			dq.Running = dq.Running[:len(dq.Running)-1]
			return
		}
	}
}

// TODO
func (y *YtDlp) saveQueueState() {}

func (y *YtDlp) download(ctx context.Context, download Download) {
	var cmdCtx context.Context
	cmd := exec.CommandContext(cmdCtx, y.Bin, download.Url, "-P", download.ID)
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
		runtime.EventsEmit(y.ctx, string(events.DownloadQueueDone), download.ID)
		y.removeFromQueue(download.ID)
		return
	}
}
