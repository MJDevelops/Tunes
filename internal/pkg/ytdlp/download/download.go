package download

import (
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/config"
)

type Download struct {
	ID       int
	Url      string
	Progress int
}

type DownloadQueue struct {
	Downloads []Download
	ctx       context.Context
}

// Threads to be used for downloads
var MaxThreads int

func init() {
	MaxThreads = config.GetMaxThreads()
	if MaxThreads <= 0 {
		MaxThreads = 5
	}
}

func (dq *DownloadQueue) SetContext(ctx context.Context) {
	dq.ctx = ctx
}

func (dq *DownloadQueue) AddToQueue(download Download) {
	dq.Downloads = append(dq.Downloads, download)
}
