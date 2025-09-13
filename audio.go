package main

import "github.com/mjdevelops/tunes/internal/pkg/audio"

// Creates a audio file given the path and returns it
func (a *App) NewAudioFile(path string) (audio.AudioFile, error) {
	return audio.NewAudioFile(path)
}
