package db

import (
	"github.com/mjdevelops/tunes/internal/pkg/db/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Track{}, &models.Artist{}, &models.Playlist{}, &models.Album{}, &models.Download{})
}
