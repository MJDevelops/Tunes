package ytdlp

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"slices"
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
	onStart    func()
}

type Queue struct {
	queue      chan Download
	ctx        context.Context
	cancel     context.CancelFunc
	once       sync.Once
	wg         sync.WaitGroup
	workers    uint
	started    atomic.Uint32
	waiting    []Download
	waitingMu  sync.Mutex
	onShutdown func([]Download)
}

func (y *YtDlp) NewDownload(url string, options ...string) Download {
	download := Download{}
	download.ID = uuid.NewString()
	download.Options = append(options, url, "--progress", "--newline", "--progress-template", "'%(progress)j'", "-q")
	download.executable = y.path
	return download
}

func NewQueue(workers uint, downloads ...Download) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	dq := &Queue{}

	dq.ctx = ctx
	dq.cancel = cancel
	dq.workers = workers
	dq.queue = make(chan Download, workers)

	for _, d := range downloads {
		dq.Enqueue(d)
	}

	return dq
}

func (d *Download) Start() (err <-chan error, cancel func()) {
	cmd := exec.CommandContext(context.Background(), d.executable, d.Options...)
	ch := make(chan error)

	if d.onStart != nil {
		d.onStart()
	}

	go func() {
		parsed := ProgressFormat{}
		r, _ := cmd.StdoutPipe()
		cmd.Start()
		if d.onProgress != nil {
			s := bufio.NewScanner(r)
			for s.Scan() {
				line := s.Bytes()
				if err := json.Unmarshal(line[1:len(line)-1], &parsed); err == nil {
					d.onProgress(parsed)
				}
			}
		}
		err := cmd.Wait()

		if d.onFinished != nil {
			d.onFinished()
		}

		ch <- err
	}()

	return ch, func() {
		cmd.Cancel()
	}
}

func (dq *Queue) OnShutdown(f func([]Download)) *Queue {
	dq.onShutdown = f
	return dq
}

func (dq *Queue) addWaiting(download Download) {
	dq.waitingMu.Lock()
	defer dq.waitingMu.Unlock()
	dq.waiting = append(dq.waiting, download)
}

func (dq *Queue) removeWaiting(id string) {
	dq.waitingMu.Lock()
	defer dq.waitingMu.Unlock()
	for i := range dq.waiting {
		if dq.waiting[i].ID == id {
			dq.waiting = slices.Delete(dq.waiting, i, i+1)
		}
	}
}

func (dq *Queue) Enqueue(download Download) {
	dq.addWaiting(download)
	go func() {
		dq.queue <- download
		dq.removeWaiting(download.ID)
	}()
}

func (dq *Queue) IsRunning() bool {
	return dq.started.Load() > 0
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
						dq.started.Add(1)
						dq.download(down)
						dq.started.Add(^uint32(0))
					}
				}
			})
		}
	})
}

func (dq *Queue) Stop() {
	dq.cancel()
	dq.wg.Wait()
	if dq.onShutdown != nil {
		dq.onShutdown(dq.waiting)
	}
}

func (dq *Queue) download(download Download) {
	dq.wg.Add(1)
	defer dq.wg.Done()
	err, cancel := download.Start()

	select {
	case <-dq.ctx.Done():
		cancel()
		os.RemoveAll(download.ID)
		dq.addWaiting(download)
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

func (d *Download) OnStart(start func()) *Download {
	d.onStart = start
	return d
}

func (d *Download) OnProgress(f func(ProgressFormat)) *Download {
	d.onProgress = f
	return d
}
