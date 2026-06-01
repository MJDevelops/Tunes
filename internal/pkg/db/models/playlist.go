package models

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Title  string  `gorm:"unique"`
	Tracks []Track `gorm:"many2many:playlist_tracks;"`
}
