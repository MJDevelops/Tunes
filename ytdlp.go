package main

import "github.com/mjdevelops/tunes/internal/pkg/ytdlp"

func (a *App) AddToDownloadQueue(download *ytdlp.Download) string {
	return a.YtDlp.AddToQueue(a.ctx, download)
}
