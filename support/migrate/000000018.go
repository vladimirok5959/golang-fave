package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000018(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(ctx, `UPDATE notify_mail SET status = 0, error = 'SMTP server is not configured' WHERE status = 2;`); err != nil {
		return err
	}

	return nil
}
