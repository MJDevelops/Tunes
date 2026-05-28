package services

import (
	"context"

	"github.com/mjdevelops/tunes/internal/pkg/audio"
	"github.com/mjdevelops/tunes/internal/pkg/db/models"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
)

type AudioService struct {
	ctx context.Context
	db  *gorm.DB
	as  *audio.AudioSink
	app *application.App
}

type AlbumWithoutTracks struct {
	ID      uint
	Title   string
	Artists []*models.Artist
}

type PlaylistWithoutTracks struct {
	ID    uint
	Title string
}

func NewAudioService(db *gorm.DB) *AudioService {
	ad := &AudioService{}
	ad.db = db
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

func (s *AudioService) GetTrack(trackId int64) (models.Track, error) {
	track, err := gorm.G[models.Track](s.db).Where("id = ?", trackId).First(s.ctx)
	if err != nil {
		return models.Track{}, err
	}
	return track, nil
}

func (s *AudioService) GetAlbumTracks(albumId int64) ([]models.Track, error) {
	album, err := gorm.G[models.Album](s.db).Where("id = ?", albumId).First(s.ctx)
	if err != nil {
		return nil, err
	}
	return album.Tracks, nil
}

func (s *AudioService) GetPlaylist(playlistId int64) (models.Playlist, error) {
	playlist, err := gorm.G[models.Playlist](s.db).Where("id = ?", playlistId).First(s.ctx)
	if err != nil {
		return models.Playlist{}, err
	}
	return playlist, nil
}

func (s *AudioService) GetPlaylists() ([]PlaylistWithoutTracks, error) {
	playlists := []PlaylistWithoutTracks{}
	err := s.db.Model(&models.Playlist{}).Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *AudioService) GetAlbums() ([]AlbumWithoutTracks, error) {
	albums := []AlbumWithoutTracks{}
	err := s.db.Model(&models.Album{}).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AudioService) Play(trackId int64, vol float64) error {
	track, err := gorm.G[models.Track](s.db).Where("id = ?", trackId).First(s.ctx)
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
