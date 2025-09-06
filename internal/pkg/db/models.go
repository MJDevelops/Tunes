package db

import (
	"database/sql"
)

type Download struct {
	ID         string `gorm:"primaryKey"`
	Url        string
	Done       bool
	FinishedAt sql.NullTime
}
