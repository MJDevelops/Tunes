package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/mjdevelops/tunes/internal/pkg/db/models"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
)

type AlbumWithoutTracks struct {
	ID      uint
	Title   string
	Artists []*models.Artist
}

type PlaylistWithoutTracks struct {
	ID    uint
	Title string
}

type DbService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewDbService(db *gorm.DB) *DbService {
	return &DbService{db: db}
}

func (s *DbService) ServiceStartup(ctx context.Context, option application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *DbService) GetTrack(trackId uint) (models.Track, error) {
	track, err := gorm.G[models.Track](s.db).Where("id = ?", trackId).First(s.ctx)
	if err != nil {
		return models.Track{}, err
	}
	return track, nil
}

func (s *DbService) GetAlbumTracks(albumId uint) ([]models.Track, error) {
	album, err := gorm.G[models.Album](s.db).Where("id = ?", albumId).First(s.ctx)
	if err != nil {
		return nil, err
	}
	return album.Tracks, nil
}

func (s *DbService) GetPlaylist(playlistId uint) (models.Playlist, error) {
	playlist, err := gorm.G[models.Playlist](s.db).Where("id = ?", playlistId).First(s.ctx)
	if err != nil {
		return models.Playlist{}, err
	}
	return playlist, nil
}

func (s *DbService) GetPlaylists() ([]PlaylistWithoutTracks, error) {
	playlists := []PlaylistWithoutTracks{}
	err := s.db.Model(&models.Playlist{}).Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *DbService) GetAlbums() ([]AlbumWithoutTracks, error) {
	albums := []AlbumWithoutTracks{}
	err := s.db.Model(&models.Album{}).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *DbService) SetDownloadFinishedAt(downloadId string, finishedAt time.Time) error {
	_, err := gorm.G[models.Download](s.db).Where("id = ?", downloadId).Update(s.ctx, "finished_at", sql.NullTime{Time: finishedAt, Valid: true})
	return err
}

func (s *DbService) GetDownload(id string) (models.Download, error) {
	download, err := gorm.G[models.Download](s.db).Where("id = ?", id).First(s.ctx)
	if err != nil {
		return models.Download{}, err
	}
	return download, nil
}

func (s *DbService) CreateDownload(download *models.Download) error {
	return gorm.G[models.Download](s.db).Create(s.ctx, download)
}

func (s *DbService) UpdateDownloadFinishedAt(id string, value time.Time) error {
	_, err := gorm.G[models.Download](s.db).Where("id = ?", id).Update(s.ctx, "finished_at", sql.NullTime{Time: value, Valid: true})
	return err
}

func (s *DbService) PendingDownloads() ([]models.Download, error) {
	downloads, err := gorm.G[models.Download](s.db).Where("finished_at IS NULL").Find(s.ctx)
	if err != nil {
		return nil, err
	}

	return downloads, nil
}
