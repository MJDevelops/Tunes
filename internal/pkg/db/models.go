package db

import (
	"database/sql"
)

type Download struct {
	ID         string `gorm:"primaryKey"`
	Url        string
	FinishedAt sql.NullTime
}
