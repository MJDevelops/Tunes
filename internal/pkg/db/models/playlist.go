package models

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Title  string
	Tracks []Track
}
