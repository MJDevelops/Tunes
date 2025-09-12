package ytdlp

import (
	"encoding/json"
	"errors"
	"os/exec"
)

type Thumbnail struct {
	Url        string      `json:"url"`
	Height     json.Number `json:"height"`
	Width      json.Number `json:"width"`
	Resolution string      `json:"resolution"`
}

func (y *YtDlp) GetThumbnails(url string) ([]Thumbnail, error) {
	var thJson struct {
		Thumbnails []Thumbnail `json:"thumbnails"`
	}

	cmd := exec.Command(y.Bin, url, "--dump-json", "-q")
	oBytes, _ := cmd.Output()

	if err := json.Unmarshal(oBytes, &thJson); err != nil {
		return nil, errors.New("couldn't parse json")
	}

	return thJson.Thumbnails, nil
}

func (y *YtDlp) GetHighDefinitionThumbnail(url string) (string, error) {
	thumbnails, err := y.GetThumbnails(url)

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
