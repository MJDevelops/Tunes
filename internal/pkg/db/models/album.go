package models

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Title   string
	Artists []*Artist `gorm:"many2many:album_artists;"`
	Tracks  []Track
}
