package db

import (
	"context"
	"database/sql"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type DB struct {
	conn *sql.DB
	ctx  context.Context
}

func NewDB() *DB {
	db, _ := sql.Open("duckdb", "")

	return &DB{
		conn: db,
	}
}

func (db *DB) SetContext(ctx context.Context) {
	db.ctx = ctx
}

func (db *DB) Close() {
	db.conn.Close()
}
