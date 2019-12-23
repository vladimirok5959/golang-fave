package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000010(ctx context.Context, db *sqlw.DB, host string) error {
	// Changes
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD COLUMN parent_id INT(11) DEFAULT NULL AFTER id;`); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD KEY FK_shop_products_parent_id (parent_id);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_parent_id
		FOREIGN KEY (parent_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	return nil
}
