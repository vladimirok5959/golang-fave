package sqlw

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"time"
)

type Tx = sql.Tx
type Rows = sql.Rows

type DB struct {
	db *sql.DB
}

var ErrNoRows = sql.ErrNoRows

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (this *DB) Close() error {
	return this.db.Close()
}

func (this *DB) Ping() error {
	return this.db.Ping()
}

func (this *DB) SetConnMaxLifetime(d time.Duration) {
	this.db.SetConnMaxLifetime(d)
}

func (this *DB) SetMaxIdleConns(n int) {
	this.db.SetMaxIdleConns(n)
}

func (this *DB) SetMaxOpenConns(n int) {
	this.db.SetMaxOpenConns(n)
}

func (this *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return this.db.QueryRow(query, args...)
}

func (this *DB) Begin() (*sql.Tx, error) {
	return this.db.Begin()
}

func (this *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return this.db.Query(query, args...)
}

func (this *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return this.db.Exec(query, args...)
}
