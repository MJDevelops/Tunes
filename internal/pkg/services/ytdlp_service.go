package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mjdevelops/tunes/internal/pkg/config"
	"github.com/mjdevelops/tunes/internal/pkg/exec/ytdlp"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type YtDlpService struct {
	ctx    context.Context
	ytdlp  *ytdlp.YtDlp
	config *config.Application
}

func NewYtDlpService(binPath string, config *config.Application) *YtDlpService {
	return &YtDlpService{
		ctx:    nil,
		ytdlp:  ytdlp.NewYtDlp(binPath),
		config: config,
	}
}

func (s *YtDlpService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	if err := s.ytdlp.GetLatest(); err != nil {
		return err
	}

	s.config.Lock()
	defer s.config.Unlock()

	s.config.YtDlp.Path = s.ytdlp.Path()
	s.config.YtDlp.Release = s.ytdlp.Release
	s.config.Write()

	return nil
}

func (s *YtDlpService) GetThumbnails(url string) ([]ytdlp.Thumbnail, error) {
	var thJson struct {
		Thumbnails []ytdlp.Thumbnail `json:"thumbnails"`
	}

	cmd := s.ytdlp.CreateCommandQuiet(url, "--dump-json")
	output, _ := cmd.Output()

	if err := json.Unmarshal(output, &thJson); err != nil {
		return nil, errors.New("couldn't parse json")
	}

	return thJson.Thumbnails, nil
}

func (s *YtDlpService) GetHighDefinitionThumbnail(url string) (string, error) {
	thumbnails, err := s.GetThumbnails(url)

	if err != nil {
		return "", err
	}

	for _, thumbnail := range thumbnails {
		if thumbnail.Resolution == "1920x1080" {
			return thumbnail.Url, nil
		}
	}

	return "", nil
}
