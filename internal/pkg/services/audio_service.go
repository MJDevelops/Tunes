package services

import (
	"context"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AudioService struct {
	ctx context.Context
	db  *DbService
	as  *audio.AudioSink
	app *application.App
}

func NewAudioService(app *application.App, dbService *DbService) *AudioService {
	ad := &AudioService{}
	ad.db = dbService
	ad.as = audio.NewAudioSink()
	ad.app = app

	audio.RegisterDecoder(&audio.TagDecoder{}, ".flac", ".ogg", ".mp3")
	audio.RegisterDecoder(&audio.WavDecoder{}, ".wav")

	return ad
}

func (s *AudioService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *AudioService) Play(trackId uint, volume float64) error {
	if s.as.IsPlaying() {
		s.as.Stop()
		// Wait for the audio sink to stop before starting a new track
		for s.as.IsPlaying() {
			time.Sleep(100 * time.Millisecond)
		}
	}

	track, err := s.db.GetTrack(trackId)
	if err != nil {
		return err
	}

	err = s.as.Init(track.Path)
	if err != nil {
		return err
	}

	err = s.as.Play(volume)
	if err != nil {
		return err
	}

	go func() {
		for {
			if s.as.IsPlaying() && !s.as.IsPaused() {
				pos := int(s.as.Position().Seconds())
				s.app.Event.EmitEvent(&application.CustomEvent{Name: "tunes:track:progress", Data: pos})
			} else if !s.as.IsPlaying() {
				s.app.Event.EmitEvent(&application.CustomEvent{Name: "tunes:track:finished"})
				return
			}
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (s *AudioService) Pause() {
	s.as.SetPlayback(false)
}

func (s *AudioService) Resume() {
	s.as.SetPlayback(true)
}

func (s *AudioService) Stop() {
	s.as.Stop()
}

func (s *AudioService) Seek(seconds int) {
	s.as.Seek(time.Duration(seconds) * time.Second)
}

func (s *AudioService) SetVolume(volume float64) {
	s.as.SetVolume(volume)
}

func (s *AudioService) Mute() {
	s.as.SetMute(true)
}

func (s *AudioService) Unmute() {
	s.as.SetMute(false)
}
