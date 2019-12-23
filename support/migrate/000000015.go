package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000015(ctx context.Context, db *sqlw.DB, host string) error {
	// Table: shop_orders
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_orders (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			create_datetime datetime NOT NULL COMMENT 'Create date/time',
			update_datetime datetime NOT NULL COMMENT 'Update date/time',
			currency_id int(11) NOT NULL COMMENT 'Currency ID',
			currency_name varchar(255) NOT NULL COMMENT 'Currency name',
			currency_coefficient float(8,4) NOT NULL DEFAULT '1.0000' COMMENT 'Currency coefficient',
			currency_code varchar(10) NOT NULL COMMENT 'Currency code',
			currency_symbol varchar(5) NOT NULL COMMENT 'Currency symbol',
			client_last_name varchar(64) NOT NULL COMMENT 'Client last name',
			client_first_name varchar(64) NOT NULL COMMENT 'Client first name',
			client_middle_name varchar(64) NOT NULL DEFAULT '' COMMENT 'Client middle name',
			client_phone varchar(20) NOT NULL DEFAULT '' COMMENT 'Client phone',
			client_email varchar(64) NOT NULL COMMENT 'Client email',
			client_delivery_comment text NOT NULL COMMENT 'Client delivery comment',
			client_order_comment text NOT NULL COMMENT 'Client order comment',
			status int(1) NOT NULL COMMENT 'new/confirmed/inprogress/canceled/completed',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_order_products
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_order_products (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			order_id int(11) NOT NULL COMMENT 'Order ID',
			product_id int(11) NOT NULL COMMENT 'Product ID',
			price float(8,2) NOT NULL COMMENT 'Product price',
			quantity int(11) NOT NULL COMMENT 'Quantity',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(ctx, `ALTER TABLE shop_orders ADD KEY FK_shop_orders_currency_id (currency_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_order_products ADD UNIQUE KEY order_product (order_id,product_id) USING BTREE;`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_order_products ADD KEY FK_shop_order_products_order_id (order_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_order_products ADD KEY FK_shop_order_products_product_id (product_id);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_orders ADD CONSTRAINT FK_shop_orders_currency_id
		FOREIGN KEY (currency_id) REFERENCES shop_currencies (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_order_products ADD CONSTRAINT FK_shop_order_products_order_id
		FOREIGN KEY (order_id) REFERENCES shop_orders (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_order_products ADD CONSTRAINT FK_shop_order_products_product_id
		FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	return nil
}
