package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name   string
	Albums []*Album `gorm:"many2many:album_artists;"`
	Tracks []Track  `gorm:"many2many:artist_tracks;"`
}
