package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000011(db *sqlw.DB, host string) error {
	if _, err := db.Exec(`ALTER TABLE shop_product_images ADD id INT NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;`); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE shop_product_images ADD COLUMN ord INT(11) NOT NULL DEFAULT 0 AFTER filename;`); err != nil {
		return err
	}

	if _, err := db.Exec(`UPDATE shop_product_images SET ord = id;`); err != nil {
		return err
	}

	return nil
}
