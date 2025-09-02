package events

type Event string

const (
	DownloadQueueStarted Event = "tunes:dqueue:started"
	DownloadQueueDone    Event = "tunes:dqueue:done"
	DownloadInterrupt    Event = "tunes:dqueue:downloadInterrupt"
	DownloadFinished     Event = "tunes:dqueue:downloadFinished"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{DownloadQueueStarted, "QUEUE_STARTED"},
	{DownloadQueueDone, "QUEUE_DONE"},
	{DownloadInterrupt, "DOWNLOAD_INTERRUPT"},
	{DownloadFinished, "DOWNLOAD_FINISHED"},
}
