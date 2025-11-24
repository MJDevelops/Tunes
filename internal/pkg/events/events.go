package events

import (
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type DownloadProgress struct {
	ID   string
	Data ytdlp.ProgressFormat
}

func RegisterWailsEvents() {
	application.RegisterEvent[application.Void]("tunes:dqueue:started")
	application.RegisterEvent[application.Void]("tunes:dqueue:done")
	application.RegisterEvent[string]("tunes:dl:started")
	application.RegisterEvent[DownloadProgress]("tunes:dl:progress")
	application.RegisterEvent[string]("tunes:dl:finished")
}
