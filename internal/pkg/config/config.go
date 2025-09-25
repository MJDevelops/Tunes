package config

import (
	"encoding/json"
	"os"
)

type YtDlp struct {
	Release string `json:"release"`
	Path    string `json:"path"`
}

type Executables struct {
	YtDlp `json:"ytdlp"`
}

type Options struct {
	MaxThreads uint `json:"maxThreads"`
}

type ApplicationConfig struct {
	Executables `json:"executables"`
	Options     `json:"options"`
	configPath  string
}

func LoadApplicationConfig(path string) (ApplicationConfig, error) {
	config := ApplicationConfig{}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0660)
	if err != nil {
		return config, err
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&config)
	config.configPath = path
	return config, nil
}

func (c *ApplicationConfig) Write() error {
	f, err := os.OpenFile(c.configPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(c)
}
