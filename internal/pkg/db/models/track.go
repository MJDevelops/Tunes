package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	Path    string    `gorm:"type:text"`
	Artists []*Artist `gorm:"many2many:artist_tracks;"`
	AlbumID uint
}
