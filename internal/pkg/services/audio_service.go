package services

import (
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AudioService struct {
	ctx context.Context
	db  *DbService
	as  *audio.AudioSink
	app *application.App
}

func NewAudioService(dbService *DbService) *AudioService {
	ad := &AudioService{}
	ad.db = dbService
	ad.as = audio.NewAudioSink()

	audio.RegisterDecoder(&audio.TagDecoder{}, ".flac", ".ogg", ".mp3")
	audio.RegisterDecoder(&audio.WavDecoder{}, ".wav")

	return ad
}

func (s *AudioService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	s.app = application.Get()
	return nil
}

func (s *AudioService) Play(trackId int64, vol float64) error {
	track, err := s.db.GetTrack(trackId)
	if err != nil {
		return err
	}

	err = s.as.Init(track.Path)
	if err != nil {
		return err
	}

	ch := s.as.Play(vol)
	go func() {
		for {
			pos, ok := <-ch
			if !ok {
				return
			}
			s.app.Event.EmitEvent(&application.CustomEvent{Name: "tunes:track:progress", Data: pos})
		}
	}()

	return nil
}
