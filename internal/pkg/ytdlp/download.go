package ytdlp

import (
	"github.com/mjdevelops/tunes/internal/pkg/config"
)

type Download struct {
	ID       int
	Url      string
	Progress int
}

type DownloadQueue struct {
	Downloads []Download
}

// Threads to be used for downloads
var MaxThreads int

func init() {
	if MaxThreads = config.GetMaxThreads(); MaxThreads <= 0 {
		MaxThreads = 5
	}
}

func (y *YtDlp) AddToQueue(download Download) {
	y.DownloadQueue.Downloads = append(y.DownloadQueue.Downloads, download)
}
