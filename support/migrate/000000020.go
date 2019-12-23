package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000020(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD COLUMN price_promo float(8,2) NOT NULL DEFAULT '0.00' AFTER price_old;`); err != nil {
		return err
	}

	return nil
}
