// Global event pool

package events

type Event string

const (
	DownloadQueueStarted Event = "tunes:dqueue:started"
	DownloadQueueDone    Event = "tunes:dqueue:done"
	DownloadStarted      Event = "tunes:dl:downloadStarted"
	DownloadInterrupt    Event = "tunes:dl:downloadInterrupt"
	DownloadFinished     Event = "tunes:dl:downloadFinished"
	DownloadProgress     Event = "tunes:dl:downloadProgress"
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
	{DownloadProgress, "DOWNLOAD_PROGRESS"},
	{TrackProgress, "TRACK_PROGRESS"},
}
