package main

import "github.com/mjdevelops/tunes/internal/pkg/audio"

// Creates a audio file given the path and returns it
func (a *App) NewAudioFile(path string) (audio.AudioFile, error) {
	ad, err := audio.NewAudioFile(path)
	if err != nil {
		return audio.AudioFile{}, err
	}
	return ad, nil
}
