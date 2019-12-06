package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000019(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN custom1 varchar(2048) NOT NULL DEFAULT '' AFTER active;`); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN custom2 varchar(2048) NOT NULL DEFAULT '' AFTER custom1;`); err != nil {
		return err
	}

	return nil
}
