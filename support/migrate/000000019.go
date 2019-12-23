package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000019(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD COLUMN custom1 varchar(2048) NOT NULL DEFAULT '' AFTER active;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD COLUMN custom2 varchar(2048) NOT NULL DEFAULT '' AFTER custom1;`); err != nil {
		return err
	}

	return nil
}
