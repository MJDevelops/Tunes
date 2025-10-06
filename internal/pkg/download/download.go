package download

import (
	"bufio"
	"context"
	"encoding/json"
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
	onProgress func(ProgressFormat)
	onFinished func()
	command    *exec.Cmd
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

func NewDownload(executable string, url string, options ...string) Download {
	download := Download{}

	opts := append(options, url, "--progress", "--newline", "--progress-template", "'%(progress)j'", "-q")
	download.command = exec.CommandContext(context.Background(), executable, opts...)

	return download
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
	ch := make(chan error)

	go func() {
		parsed := ProgressFormat{}
		r, _ := download.command.StdoutPipe()
		download.command.Start()
		s := bufio.NewScanner(r)
		for s.Scan() {
			line := s.Bytes()
			if err := json.Unmarshal(line[1:len(line)-1], &parsed); err == nil {
				download.onProgress(parsed)
			}
		}
		ch <- download.command.Wait()
	}()

	select {
	case <-dq.ctx.Done():
		download.command.Cancel()
		os.RemoveAll(download.ID)
	case <-ch:
		download.onFinished()
	}
}

func (d *Download) OnFinished(f func()) *Download {
	d.onFinished = f
	return d
}

func (d *Download) OnProgress(f func(ProgressFormat)) *Download {
	d.onProgress = f
	return d
}
