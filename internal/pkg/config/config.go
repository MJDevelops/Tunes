package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type YtDlp struct {
	Release string `json:"release"`
	Path    string `json:"path"`
}

type Ffmpeg struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}

type Executables struct {
	YtDlp  `json:"ytdlp"`
	Ffmpeg `json:"ffmpeg"`
}

type Options struct {
	MaxThreads uint `json:"maxThreads"`
}

type Application struct {
	Executables `json:"executables"`
	Options     `json:"options"`
	configPath  string
}

func LoadApplicationConfig(path string) (Application, error) {
	config := Application{}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0660)
	if err != nil {
		return config, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&config); !errors.Is(err, io.EOF) && err != nil {
		return config, err
	}
	config.configPath = path

	return config, nil
}

func (c *Application) Write() error {
	f, err := os.OpenFile(c.configPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	pretty, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	f.Write(pretty)

	return nil
}
