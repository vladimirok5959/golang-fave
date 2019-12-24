package sqlw

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"time"

	"golang-fave/engine/consts"
)

type Tx struct {
	tx *sql.Tx
	s  time.Time
}

func (this *Tx) Commit() error {
	if consts.ParamDebug {
		err := this.tx.Commit()
		log("[TX] TRANSACTION END (Commit)", this.s, err, true)
		return err
	}
	return this.tx.Commit()
}

func (this *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.tx.ExecContext(ctx, query, args...)
		log("[TX] "+query, s, e, true)
		return r, e
	}
	return this.tx.ExecContext(ctx, query, args...)
}

func (this *Tx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.tx.QueryContext(ctx, query, args...)
		log("[TX] "+query, s, e, true)
		return r, e
	}
	return this.tx.QueryContext(ctx, query, args...)
}

func (this *Tx) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if consts.ParamDebug {
		s := time.Now()
		r := this.tx.QueryRowContext(ctx, query, args...)
		log("[TX] "+query, s, nil, true)
		return r
	}
	return this.tx.QueryRowContext(ctx, query, args...)
}

func (this *Tx) Rollback() error {
	if consts.ParamDebug {
		err := this.tx.Rollback()
		log("[TX] TRANSACTION END (Rollback)", this.s, nil, true)
		return err
	}
	return this.tx.Rollback()
}
