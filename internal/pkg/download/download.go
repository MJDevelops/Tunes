package download

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type ProgressFormat struct {
	ETA             uint    `json:"eta"`
	Speed           float32 `json:"speed"`
	Elapsed         float32 `json:"elapsed"`
	DownloadedBytes uint    `json:"downloaded_bytes"`
	TotalBytes      uint    `json:"total_bytes"`
}

type Download struct {
	ID         string
	Url        string
	OnProgress func(ProgressFormat)
	OnFinished func(*Download)
	command    string
}

type DownloadQueue struct {
	running []Download
	waiting []Download
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.RWMutex
	once    sync.Once
	wg      sync.WaitGroup
	workers uint
}

func NewDownloadQueue(workers uint, downloads ...Download) *DownloadQueue {
	ctx, cancel := context.WithCancel(context.Background())
	dq := &DownloadQueue{}

	dq.ctx = ctx
	dq.cancel = cancel
	dq.workers = workers

	return dq
}

func (dq *DownloadQueue) Running() []Download {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	return dq.running
}

func (dq *DownloadQueue) Waiting() []Download {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	return dq.waiting
}

// Adds download to queue and returns the corresponding ID
func (dq *DownloadQueue) AddToQueue(download Download) string {
	id := uuid.NewString()
	download.ID = id

	dq.mu.Lock()
	defer dq.mu.Unlock()

	dq.waiting = append(dq.waiting, download)
	return id
}

func (dq *DownloadQueue) IsRunning() bool {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	return len(dq.running) > 0 || len(dq.waiting) > 0
}

func (dq *DownloadQueue) Start() {
	dq.once.Do(func() {
		for range dq.workers {
			dq.wg.Add(1)
			go func() {
				defer dq.wg.Done()
				for {
					select {
					case <-dq.ctx.Done():
						return
					default:
						if dq.isWaiting() {
							dq.mu.RLock()
							newDown := dq.waiting[0]
							dq.mu.RUnlock()

							dq.download(newDown)

							dq.mu.Lock()
							dq.running = append(dq.running, newDown)
							dq.waiting = slices.Delete(dq.waiting, 0, 1)
							dq.mu.Unlock()
						}
					}
				}
			}()
		}
	})
}

func (dq *DownloadQueue) isWaiting() bool {
	dq.mu.RLock()
	defer dq.mu.RUnlock()

	return len(dq.waiting) > 0
}

func (dq *DownloadQueue) removeFromQueue(id string) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for i, d := range dq.running {
		if d.ID == id {
			dq.running = slices.Delete(dq.running, i, i+1)
			return
		}
	}
}

func (dq *DownloadQueue) Stop() {
	dq.cancel()
	dq.wg.Wait()
}

func (dq *DownloadQueue) download(download Download) {
	dq.wg.Add(1)
	defer dq.wg.Done()
	cmd := exec.Command(download.command, download.Url, "--progress", "--newline", "--progress-template", "'%(progress)j'", "-q")
	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	select {
	case <-dq.ctx.Done():
		cmd.Cancel()
		os.RemoveAll(download.ID)
		return
	case <-ch:
		dq.removeFromQueue(download.ID)
		return
	}
}
