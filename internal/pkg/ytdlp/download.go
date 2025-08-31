package ytdlp

import (
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
}

// Threads to be used for downloads
var MaxThreads int

func init() {
	if MaxThreads = config.GetMaxThreads(); MaxThreads <= 0 {
		MaxThreads = 5
	}
}

// Adds download to queue and returns the corresponding ID
func (y *YtDlp) AddToQueue(download Download) string {
	id := uuid.NewString()
	download.ID = id
	y.DownloadQueue.Downloads = append(y.DownloadQueue.Downloads, download)
	return id
}
