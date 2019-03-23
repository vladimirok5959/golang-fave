package sqlw

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"time"
)

type Tx struct {
	tx *sql.Tx
	s  time.Time
}

func (this *Tx) Rollback() error {
	err := this.tx.Rollback()
	log("[TX] TRANSACTION END (Rollback)", this.s, true)
	return err
}

func (this *Tx) Commit() error {
	err := this.tx.Commit()
	log("[TX] TRANSACTION END (Commit)", this.s, true)
	return err
}

func (this *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	s := time.Now()
	r, e := this.tx.Exec(query, args...)
	log("[TX] "+query, s, true)
	return r, e
}

func (this *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	s := time.Now()
	r := this.tx.QueryRow(query, args...)
	log("[TX] "+query, s, true)
	return r
}
