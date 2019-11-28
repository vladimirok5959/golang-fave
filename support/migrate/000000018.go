package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000018(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`UPDATE notify_mail SET status = 0, error = 'SMTP server is not configured' WHERE status = 2;`); err != nil {
		return err
	}

	return nil
}
