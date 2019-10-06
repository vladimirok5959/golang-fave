package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000013(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN gname VARCHAR(255) NOT NULL DEFAULT '' AFTER price;`); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE shop_products ALTER gname DROP DEFAULT;`); err != nil {
		return err
	}

	return nil
}
