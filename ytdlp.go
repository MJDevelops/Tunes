package main

import (
	"encoding/json"
	"errors"

	"github.com/mjdevelops/tunes/internal/pkg/ytdlp"
)

func (a *App) GetThumbnails(url string) ([]ytdlp.Thumbnail, error) {
	var thJson struct {
		Thumbnails []ytdlp.Thumbnail `json:"thumbnails"`
	}

	cmd := a.YtDlp.CreateCommandQuiet(url, "--dump-json")
	output, _ := cmd.Output()

	if err := json.Unmarshal(output, &thJson); err != nil {
		return nil, errors.New("couldn't parse json")
	}

	return thJson.Thumbnails, nil
}

func (a *App) GetHighDefinitionThumbnail(url string) (string, error) {
	thumbnails, err := a.GetThumbnails(url)

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
