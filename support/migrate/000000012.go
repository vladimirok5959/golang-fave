package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000012(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD KEY name (name);`); err != nil {
		return err
	}

	return nil
}
