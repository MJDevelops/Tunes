package download

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"

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
	queue      chan Download
	ctx        context.Context
	cancel     context.CancelFunc
	once       sync.Once
	wg         sync.WaitGroup
	workers    uint
	started    uint32
	onShutdown func(downloads <-chan Download)
}

func NewDownloadQueue(workers uint, downloads ...Download) *DownloadQueue {
	ctx, cancel := context.WithCancel(context.Background())
	dq := &DownloadQueue{}

	dq.ctx = ctx
	dq.cancel = cancel
	dq.workers = workers
	dq.queue = make(chan Download, len(downloads))

	for _, d := range downloads {
		dq.queue <- d
	}

	return dq
}

func (dq *DownloadQueue) OnShutdown(f func(downloads <-chan Download)) *DownloadQueue {
	dq.onShutdown = f
	return dq
}

// Adds download to queue and returns the corresponding ID
func (dq *DownloadQueue) SendToQueue(download Download) string {
	id := uuid.NewString()
	download.ID = id

	go func() {
		dq.queue <- download
	}()

	return id
}

func (dq *DownloadQueue) IsRunning() bool {
	return atomic.LoadUint32(&dq.started) > 0
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
					case down := <-dq.queue:
						atomic.AddUint32(&dq.started, 1)
						dq.download(down)
						atomic.AddUint32(&dq.started, ^uint32(0))
					}
				}
			}()
		}
	})
}

func (dq *DownloadQueue) Stop() {
	dq.cancel()
	dq.wg.Wait()
	dq.onShutdown(dq.queue)
}

func (dq *DownloadQueue) download(download Download) {
	dq.wg.Add(1)
	defer dq.wg.Done()
	cmd := exec.Command(download.command, download.Url, "--progress", "--newline", "--progress-template", "'%(progress)j'")
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
		download.OnFinished(&download)
		return
	}
}
