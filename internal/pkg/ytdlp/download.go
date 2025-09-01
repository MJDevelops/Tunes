package ytdlp

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
	"github.com/mjdevelops/tunes/internal/pkg/config"
)

type Download struct {
	ID       string
	Url      string
	Progress int
}

type DownloadQueue struct {
	Downloads []Download
	mu        sync.Mutex
	once      sync.Once
}

var (
	// Threads to be used for downloads
	MaxThreads int
	dc         chan *Download
)

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
	dq.Downloads = append(dq.Downloads, *download)
	dc <- download
	return id
}

func (y *YtDlp) StartQueue(ctx context.Context) {
	y.DownloadQueue.once.Do(func() {
		throttle := make(chan int, MaxThreads)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					download := <-dc
					throttle <- 1
					go func() {
						y.startDownload(ctx, download)
						y.RemoveFromQueue(ctx, download.ID)
						<-throttle
					}()
				}
			}
		}()
	})
}

func (y *YtDlp) StopQueue() {
	y.DownloadQueue.mu.Lock()
	y.saveQueueState()
}

func (y *YtDlp) RemoveFromQueue(ctx context.Context, id string) {
	dq := &y.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for i, d := range dq.Downloads {
		if d.ID == id {
			dq.Downloads[i] = dq.Downloads[len(dq.Downloads)-1]
			dq.Downloads = dq.Downloads[:len(dq.Downloads)-1]
			break
		}
	}
}

// TODO
func (y *YtDlp) saveQueueState() {}

func (y *YtDlp) startDownload(ctx context.Context, download *Download) {
	var cmdCtx context.Context
	cmd := exec.CommandContext(cmdCtx, y.Bin, download.Url, "-P", download.ID)
	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	select {
	case <-ctx.Done():
		cmd.Cancel()
		os.RemoveAll(download.ID)
		return
	case <-ch:
		return
	}
}
