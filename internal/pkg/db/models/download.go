package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Download struct {
	gorm.Model
	ID         string
	FinishedAt sql.NullTime
	Options    string `gorm:"type:text"`
	Source     string
	Path       string
}
