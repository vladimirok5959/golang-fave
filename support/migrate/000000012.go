package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000012(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`ALTER TABLE shop_products ADD KEY name (name);`); err != nil {
		return err
	}

	return nil
}
