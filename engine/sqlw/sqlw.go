package sqlw

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"time"

	"golang-fave/consts"
)

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
	return &DB{db: db}, err
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
	if consts.ParamDebug {
		s := time.Now()
		r := this.db.QueryRow(query, args...)
		log(query, s, false)
		return r
	}
	return this.db.QueryRow(query, args...)
}

func (this *DB) Begin() (*Tx, error) {
	tx, err := this.db.Begin()
	if err != nil {
		return nil, err
	}
	if consts.ParamDebug {
		s := time.Now()
		log("[TX] TRANSACTION START", s, true)
		return &Tx{tx, s}, err
	}
	return &Tx{tx, time.Now()}, err
}

func (this *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.db.Query(query, args...)
		log(query, s, false)
		return r, e
	}
	return this.db.Query(query, args...)
}

func (this *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.db.Exec(query, args...)
		log(query, s, false)
		return r, e
	}
	return this.db.Exec(query, args...)
}

func (this *DB) Transaction(queries func(tx *Tx) error) error {
	if queries == nil {
		return errors.New("queries is not set for transaction")
	}
	tx, err := this.Begin()
	if err != nil {
		return err
	}
	err = queries(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
