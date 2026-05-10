package services

import (
	"container/list"
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/db/models"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
)

type AudioService struct {
	ctx   context.Context
	db    *gorm.DB
	elems map[int64]*list.Element
}

func NewAudioService(db *gorm.DB) *AudioService {
	ad := &AudioService{}
	ad.db = db

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

func (s *AudioService) GetPlaylistTracks(playlistId int64) ([]models.Track, error) {
	ctx := context.Background()
	playlist, err := gorm.G[models.Playlist](s.db).Where("id = ?", playlistId).First(ctx)
	if err != nil {
		return nil, err
	}
	return playlist.Tracks, nil
}
