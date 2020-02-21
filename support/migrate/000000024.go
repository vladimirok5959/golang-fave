package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000024(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(ctx, `ALTER TABLE fave_pages ADD COLUMN template varchar(255) NOT NULL DEFAULT 'page' AFTER user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `UPDATE fave_pages SET template = 'index' WHERE alias = '/' LIMIT 1;`); err != nil {
		return err
	}

	return nil
}
