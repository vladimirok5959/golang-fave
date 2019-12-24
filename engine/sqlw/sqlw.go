package sqlw

import (
	"context"
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
		if consts.ParamDebug {
			log("[CM] OPEN", time.Now(), err, true)
		}
		return nil, err
	}
	if consts.ParamDebug {
		log("[CM] OPEN", time.Now(), err, true)
	}
	return &DB{db: db}, err
}

func (this *DB) Begin(ctx context.Context) (*Tx, error) {
	tx, err := this.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		if consts.ParamDebug {
			log("[TX] TRANSACTION START", time.Now(), err, true)
		}
		return nil, err
	}
	if consts.ParamDebug {
		s := time.Now()
		log("[TX] TRANSACTION START", s, err, true)
		return &Tx{tx, s}, err
	}
	return &Tx{tx, time.Now()}, err
}

func (this *DB) Close() error {
	if consts.ParamDebug {
		err := this.db.Close()
		log("[CM] CLOSE", time.Now(), err, true)
		return err
	}
	return this.db.Close()
}

func (this *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.db.ExecContext(ctx, query, args...)
		log(query, s, e, false)
		return r, e
	}
	return this.db.ExecContext(ctx, query, args...)
}

func (this *DB) Ping(ctx context.Context) error {
	if consts.ParamDebug {
		err := this.db.PingContext(ctx)
		log("[CM] PING", time.Now(), err, true)
		return err
	}
	return this.db.PingContext(ctx)
}

func (this *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.db.QueryContext(ctx, query, args...)
		log(query, s, e, false)
		return r, e
	}
	return this.db.Query(query, args...)
}

func (this *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if consts.ParamDebug {
		s := time.Now()
		r := this.db.QueryRowContext(ctx, query, args...)
		log(query, s, nil, false)
		return r
	}
	return this.db.QueryRow(query, args...)
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

func (this *DB) Transaction(ctx context.Context, queries func(ctx context.Context, tx *Tx) error) error {
	if queries == nil {
		return errors.New("queries is not set for transaction")
	}
	tx, err := this.Begin(ctx)
	if err != nil {
		return err
	}
	if err := queries(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
