package ytdlp

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
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
	Options    []string
	executable string
	onProgress func(ProgressFormat)
	onFinished func()
}

type Queue struct {
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
	download.ID = uuid.NewString()
	download.Options = append(options, url, "--progress", "--newline", "--progress-template", "'%(progress)j'", "-q")
	download.executable = executable
	return download
}

func NewQueue(workers uint, downloads ...Download) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	dq := &Queue{}

	dq.ctx = ctx
	dq.cancel = cancel
	dq.workers = workers
	dq.queue = make(chan Download, len(downloads))

	for _, d := range downloads {
		dq.queue <- d
	}

	return dq
}

func (d *Download) Start() (err <-chan error, cancel func()) {
	cmd := exec.CommandContext(context.Background(), d.executable, d.Options...)
	ch := make(chan error)

	go func() {
		parsed := ProgressFormat{}
		r, _ := cmd.StdoutPipe()
		cmd.Start()
		s := bufio.NewScanner(r)
		for s.Scan() {
			line := s.Bytes()
			if err := json.Unmarshal(line[1:len(line)-1], &parsed); err == nil {
				d.onProgress(parsed)
			}
		}
		err := cmd.Wait()
		d.onFinished()
		ch <- err
	}()

	return ch, func() {
		cmd.Cancel()
	}
}

func (dq *Queue) OnShutdown(f func(downloads <-chan Download)) *Queue {
	dq.onShutdown = f
	return dq
}

func (dq *Queue) SendToQueue(download Download) {
	go func() {
		dq.queue <- download
	}()
}

func (dq *Queue) IsRunning() bool {
	return atomic.LoadUint32(&dq.started) > 0
}

func (dq *Queue) Start() {
	dq.once.Do(func() {
		for range dq.workers {
			dq.wg.Go(func() {
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
			})
		}
	})
}

func (dq *Queue) Stop() {
	dq.cancel()
	dq.wg.Wait()
	dq.onShutdown(dq.queue)
}

func (dq *Queue) download(download Download) {
	dq.wg.Add(1)
	defer dq.wg.Done()
	err, cancel := download.Start()

	select {
	case <-dq.ctx.Done():
		cancel()
		os.RemoveAll(download.ID)
	case err := <-err:
		if err != nil {
			log.Printf("Error executing download with ID %s: %s", download.ID, err.Error())
		} else {
			log.Printf("Download with ID %s finished", download.ID)
		}
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
