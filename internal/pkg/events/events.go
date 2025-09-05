package events

type Event string

const (
	DownloadQueueStarted Event = "tunes:dqueue:started"
	DownloadQueueDone    Event = "tunes:dqueue:done"
	DownloadStarted      Event = "tunes:dqueue:downloadStarted"
	DownloadInterrupt    Event = "tunes:dqueue:downloadInterrupt"
	DownloadFinished     Event = "tunes:dqueue:downloadFinished"
	TrackProgress        Event = "tunes:track:progress"
)

var Events = []struct {
	Value  Event
	TSName string
}{
	{DownloadQueueStarted, "QUEUE_STARTED"},
	{DownloadQueueDone, "QUEUE_DONE"},
	{DownloadStarted, "DOWNLOAD_STARTED"},
	{DownloadInterrupt, "DOWNLOAD_INTERRUPT"},
	{DownloadFinished, "DOWNLOAD_FINISHED"},
	{TrackProgress, "TRACK_PROGRESS"},
}
