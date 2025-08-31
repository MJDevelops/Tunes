package ytdlp

import (
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
func (q *DownloadQueue) AddToQueue(download Download) string {
	id := uuid.NewString()
	download.ID = id
	q.Downloads = append(q.Downloads, download)
	return id
}

func (q *DownloadQueue) StartQueue() {
	throttle := make(chan int, 5)
	var wg sync.WaitGroup
	for _, download := range q.Downloads {
		throttle <- 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			download.start()
			q.RemoveFromQueue(download)
			<-throttle
		}()
	}
	wg.Wait()
}

func (q *DownloadQueue) RemoveFromQueue(download Download) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, d := range q.Downloads {
		if d.ID == download.ID {
			q.Downloads[i] = q.Downloads[len(q.Downloads)-1]
			q.Downloads = q.Downloads[:len(q.Downloads)-1]
			break
		}
	}
}
