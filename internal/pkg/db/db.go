package db

import (
	"database/sql"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type DB struct {
	Conn *sql.DB
}

func NewDB() (*DB, error) {
	conn, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}

	return &DB{
		Conn: conn,
	}, nil
}

func (d *DB) Close() {
	d.Conn.Close()
}
