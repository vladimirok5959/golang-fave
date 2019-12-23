package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000006(ctx context.Context, db *sqlw.DB, host string) error {
	// Table: shop_product_images
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_product_images (
			product_id int(11) NOT NULL,
			filename varchar(255) NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(ctx, `ALTER TABLE shop_product_images ADD UNIQUE KEY product_filename (product_id,filename) USING BTREE;`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_product_images ADD KEY FK_shop_product_images_product_id (product_id);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_product_images ADD CONSTRAINT FK_shop_product_images_product_id
		FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	return nil
}
