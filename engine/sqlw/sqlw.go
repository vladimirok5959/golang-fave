package sqlw

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang-fave/consts"
)

type Tx = sql.Tx
type Rows = sql.Rows

type DB struct {
	db *sql.DB
}

var ErrNoRows = sql.ErrNoRows

func (this *DB) logQuery(query string, s time.Time) {
	msg := query
	if reg, err := regexp.Compile("[\\s\\t]+"); err == nil {
		msg = strings.Trim(reg.ReplaceAllString(msg, " "), " ")
	}
	if consts.ParamDebug {
		t := time.Now().Sub(s).Seconds()
		if consts.IS_WIN {
			fmt.Fprintln(os.Stdout, "[SQL] "+msg+fmt.Sprintf(" %.3f ms", t))
		} else {
			fmt.Fprintln(os.Stdout, "\033[1;33m[SQL] "+msg+fmt.Sprintf(" %.3f ms", t)+"\033[0m")
		}
	}
}

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
	s := time.Now()
	r := this.db.QueryRow(query, args...)
	this.logQuery(query, s)
	return r
}

func (this *DB) Begin() (*sql.Tx, error) {
	return this.db.Begin()
}

func (this *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s := time.Now()
	r, e := this.db.Query(query, args...)
	this.logQuery(query, s)
	return r, e
}

func (this *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	s := time.Now()
	r, e := this.db.Exec(query, args...)
	this.logQuery(query, s)
	return r, e
}
