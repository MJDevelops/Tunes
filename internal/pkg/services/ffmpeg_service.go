package services

import (
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ffmpeg"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type FfmpegService struct {
	ctx    context.Context
	ffmpeg *ffmpeg.Ffmpeg
	config *config.Application
}

func NewFfmpegService(binPath string, config *config.Application) *FfmpegService {
	ffmpeg, _ := ffmpeg.NewFfmpeg(binPath)
	return &FfmpegService{
		ffmpeg: ffmpeg,
		ctx:    nil,
		config: config,
	}
}

func (fs *FfmpegService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	fs.ctx = ctx

	if err := fs.ffmpeg.GetLatest(); err != nil {
		return err
	}

	fs.config.Lock()
	defer fs.config.Unlock()

	fs.config.Ffmpeg.Path = fs.ffmpeg.Path()
	fs.config.Ffmpeg.Version = fs.ffmpeg.Version()
	fs.config.Write()

	return nil
}
