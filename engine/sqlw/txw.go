package sqlw

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"time"

	"golang-fave/consts"
)

type Tx struct {
	tx *sql.Tx
	s  time.Time
}

func (this *Tx) Rollback() error {
	if consts.ParamDebug {
		err := this.tx.Rollback()
		log("[TX] TRANSACTION END (Rollback)", this.s, true)
		return err
	}
	return this.tx.Rollback()
}

func (this *Tx) Commit() error {
	if consts.ParamDebug {
		err := this.tx.Commit()
		log("[TX] TRANSACTION END (Commit)", this.s, true)
		return err
	}
	return this.tx.Commit()
}

func (this *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	if consts.ParamDebug {
		s := time.Now()
		r, e := this.tx.Exec(query, args...)
		log("[TX] "+query, s, true)
		return r, e
	}
	return this.tx.Exec(query, args...)
}

func (this *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	if consts.ParamDebug {
		s := time.Now()
		r := this.tx.QueryRow(query, args...)
		log("[TX] "+query, s, true)
		return r
	}
	return this.tx.QueryRow(query, args...)
}
