package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000007(db *sqlw.DB, host string) error {
	// Changes
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN vendor VARCHAR(255) NOT NULL DEFAULT '' AFTER alias;`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN quantity INT(11) NOT NULL DEFAULT 0 AFTER vendor;`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE shop_products ADD COLUMN category INT(11) NOT NULL DEFAULT 1 AFTER quantity;`); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(`ALTER TABLE shop_products ADD KEY FK_shop_products_category (category);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(`
		ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_category
		FOREIGN KEY (category) REFERENCES shop_cats (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	// Remove default
	if _, err := db.Exec(`ALTER TABLE shop_products ALTER vendor DROP DEFAULT;`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE shop_products ALTER quantity DROP DEFAULT;`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE shop_products ALTER category DROP DEFAULT;`); err != nil {
		return err
	}

	return nil
}
