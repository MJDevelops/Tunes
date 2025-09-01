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
	once      sync.Once
}

// Threads to be used for downloads
var MaxThreads int
var exitSig = make(chan bool)
var dc chan *Download
var wg sync.WaitGroup

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

func (y *YtDlp) StartQueue() {
	y.DownloadQueue.once.Do(func() {
		throttle := make(chan int, MaxThreads)
		go func() {
			for {
				select {
				case <-exitSig:
					return
				default:
					download := <-dc
					throttle <- 1
					go func() {
						download.start()
						wg.Add(1)
						defer wg.Done()
						y.RemoveFromQueue(download.ID)
						<-throttle
					}()
				}
			}
		}()
	})
}

func (y *YtDlp) StopQueue() {
	exitSig <- true
	y.DownloadQueue.mu.Lock()
	wg.Wait()
	y.saveQueueState()
}

func (y *YtDlp) RemoveFromQueue(id string) {
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

// TODO
func (d *Download) start() {}
