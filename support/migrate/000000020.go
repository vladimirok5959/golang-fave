package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000020(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN price_promo float(8,2) NOT NULL DEFAULT '0.00' AFTER price_old;`); err != nil {
		return err
	}

	return nil
}
