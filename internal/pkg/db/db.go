package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("tunes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{
		Conn: db,
	}, nil
}
