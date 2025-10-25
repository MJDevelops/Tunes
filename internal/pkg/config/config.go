package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
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
	MaxThreads uint `json:"maxThreads,omitempty"`
}

type Application struct {
	Executables `json:"executables"`
	Options     `json:"options"`
	configPath  string
	mu          sync.Mutex
}

func LoadApplicationConfig(path string) (*Application, error) {
	config := &Application{}
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

// Lock acquires a lock on the application configuration.
// This function should be called before modifying the configuration.
func (c *Application) Lock() {
	c.mu.Lock()
}

// Unlock releases the lock on the application configuration.
// This function needs to be called after acquiring the lock with Lock.
func (c *Application) Unlock() {
	c.mu.Unlock()
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
