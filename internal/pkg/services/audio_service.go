package services

import (
	"container/list"
	"context"

	"github.com/mjdevelops/tunes/db"
	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AudioService struct {
	ctx     context.Context
	queries *db.Queries
	elems   map[int64]*list.Element
}

func NewAudioService(queries *db.Queries) *AudioService {
	ad := &AudioService{}
	ad.queries = queries

	audio.RegisterDecoder(&audio.TagDecoder{}, ".flac", ".ogg", ".mp3")
	audio.RegisterDecoder(&audio.WavDecoder{}, ".wav")

	return ad
}

func (s *AudioService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *AudioService) AddToQueue(trackId int64) error {
	return nil
}

func (s *AudioService) GetPlaylistTracks(playlistId int64) ([]db.Track, error) {
	return s.queries.GetPlaylistTracks(context.TODO(), playlistId)
}
