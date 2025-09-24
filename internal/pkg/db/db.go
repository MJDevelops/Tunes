package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	conn *gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("tunes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{
		conn: db,
	}, nil
}

func (db *DB) Migrate() {
	models := []any{
		&Download{},
	}
	db.conn.AutoMigrate(models...)
}

func (db *DB) Conn() *gorm.DB {
	return db.conn
}
