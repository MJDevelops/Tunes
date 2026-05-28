package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	Title   string
	Artists []*Artist `gorm:"many2many:artist_tracks;"`
	Path    string    `gorm:"type:text"`
	AlbumID uint
}
