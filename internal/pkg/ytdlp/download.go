package ytdlp

import (
	"sync"

	"github.com/google/uuid"
	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Download struct {
	ID       string
	Url      string
	Progress int
}

type DownloadQueue struct {
	Downloads []Download
	mu        sync.Mutex
}

// Threads to be used for downloads
var MaxThreads int

func init() {
	if MaxThreads = config.GetMaxThreads(); MaxThreads <= 0 {
		MaxThreads = 5
	}
}

func (d *Download) start() {}

// Adds download to queue and returns the corresponding ID
func (y *YtDlp) AddToQueue(download Download) string {
	dq := &y.DownloadQueue
	id := uuid.NewString()
	download.ID = id
	dq.mu.Lock()
	dq.Downloads = append(dq.Downloads, download)
	dq.mu.Unlock()
	return id
}

func (y *YtDlp) StartQueue() {
	throttle := make(chan int, 5)
	var wg sync.WaitGroup
	for _, download := range y.DownloadQueue.Downloads {
		throttle <- 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			download.start()
			y.RemoveFromQueue(download)
			<-throttle
		}()
	}
	wg.Wait()
	runtime.EventsEmit(y.ctx, "tunes:dqueue:done")
}

func (y *YtDlp) RemoveFromQueue(download Download) {
	dq := &y.DownloadQueue
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for i, d := range dq.Downloads {
		if d.ID == download.ID {
			dq.Downloads[i] = dq.Downloads[len(dq.Downloads)-1]
			dq.Downloads = dq.Downloads[:len(dq.Downloads)-1]
			break
		}
	}
}
