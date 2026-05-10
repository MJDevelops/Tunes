package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	Path string `gorm:"type:text"`
}
