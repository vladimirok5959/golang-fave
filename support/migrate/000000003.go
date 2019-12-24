package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"
)

func Migrate_000000003(ctx context.Context, db *sqlw.DB, host string) error {
	// Remove blog indexes
	if _, err := db.Exec(ctx, `DROP INDEX post_id ON blog_cat_post_rel`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `DROP INDEX category_id ON blog_cat_post_rel`); err != nil {
		return err
	}

	// Table: shop_cat_product_rel
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_cat_product_rel (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			product_id int(11) NOT NULL COMMENT 'Product id',
			category_id int(11) NOT NULL COMMENT 'Category id',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_cats
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_cats (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			user int(11) NOT NULL COMMENT 'User id',
			name varchar(255) NOT NULL COMMENT 'Category name',
			alias varchar(255) NOT NULL COMMENT 'Category alias',
			lft int(11) NOT NULL COMMENT 'For nested set model',
			rgt int(11) NOT NULL COMMENT 'For nested set model',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_currencies
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_currencies (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			name varchar(255) NOT NULL COMMENT 'Currency name',
			coefficient float(8,4) NOT NULL DEFAULT '1.0000' COMMENT 'Currency coefficient',
			code varchar(10) NOT NULL COMMENT 'Currency code',
			symbol varchar(5) NOT NULL COMMENT 'Currency symbol',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_filter_product_values
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_filter_product_values (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			product_id int(11) NOT NULL COMMENT 'Product id',
			filter_value_id int(11) NOT NULL COMMENT 'Filter value id',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_filters
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_filters (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			name varchar(255) NOT NULL COMMENT 'Filter name in CP',
			filter varchar(255) NOT NULL COMMENT 'Filter name in site',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_filters_values
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_filters_values (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			filter_id int(11) NOT NULL COMMENT 'Filter id',
			name varchar(255) NOT NULL COMMENT 'Value name',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Table: shop_products
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE shop_products (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			user int(11) NOT NULL COMMENT 'User id',
			currency int(11) NOT NULL COMMENT 'Currency id',
			price float(8,2) NOT NULL COMMENT 'Product price',
			name varchar(255) NOT NULL COMMENT 'Product name',
			alias varchar(255) NOT NULL COMMENT 'Product alias',
			briefly text NOT NULL COMMENT 'Product brief content',
			content text NOT NULL COMMENT 'Product content',
			datetime datetime NOT NULL COMMENT 'Creation date/time',
			active int(1) NOT NULL COMMENT 'Is active product or not',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Demo datas
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_cat_product_rel (id, product_id, category_id)
			VALUES
		(1, 1, 3);`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_cats (id, user, name, alias, lft, rgt)
			VALUES
		(1, 1, 'ROOT', 'ROOT', 1, 6),
		(2, 1, 'Electronics', 'electronics', 2, 5),
		(3, 1, 'Mobile phones', 'mobile-phones', 3, 4);`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_currencies (id, name, coefficient, code, symbol)
			VALUES
		(1, 'US Dollar', 1.0000, 'USD', '$');`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_filter_product_values (id, product_id, filter_value_id)
			VALUES
		(1, 1, 3),
		(2, 1, 7),
		(3, 1, 9),
		(4, 1, 10),
		(5, 1, 11);`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_filters (id, name, filter)
			VALUES
		(1, 'Mobile phones manufacturer', 'Manufacturer'),
		(2, 'Mobile phones memory', 'Memory'),
		(3, 'Mobile phones communication standard', 'Communication standard');`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_filters_values (id, filter_id, name)
			VALUES
		(1, 1, 'Apple'),
		(2, 1, 'Asus'),
		(3, 1, 'Samsung'),
		(4, 2, '16 Gb'),
		(5, 2, '32 Gb'),
		(6, 2, '64 Gb'),
		(7, 2, '128 Gb'),
		(8, 2, '256 Gb'),
		(9, 3, '4G'),
		(10, 3, '2G'),
		(11, 3, '3G');`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`INSERT INTO shop_products SET
			id = ?,
			user = ?,
			currency = ?,
			price = ?,
			name = ?,
			alias = ?,
			briefly = ?,
			content = ?,
			datetime = ?,
			active = ?
		;`,
		1,
		1,
		1,
		1000.00,
		"Samsung Galaxy S10",
		"samsung-galaxy-s10",
		"<p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent.</p>",
		"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
		utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
		1,
	); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cat_product_rel ADD UNIQUE KEY product_category (product_id,category_id) USING BTREE;`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_product_id (product_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_category_id (category_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cats ADD UNIQUE KEY alias (alias);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cats ADD KEY lft (lft), ADD KEY rgt (rgt);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_cats ADD KEY FK_shop_cats_user (user);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filter_product_values ADD UNIQUE KEY product_filter_value (product_id,filter_value_id) USING BTREE;`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_product_id (product_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_filter_value_id (filter_value_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filters ADD KEY name (name);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filters_values ADD KEY FK_shop_filters_values_filter_id (filter_id);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_filters_values ADD KEY name (name);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD UNIQUE KEY alias (alias);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD KEY FK_shop_products_user (user);`); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, `ALTER TABLE shop_products ADD KEY FK_shop_products_currency (currency);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_product_id
		FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_category_id
		FOREIGN KEY (category_id) REFERENCES shop_cats (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_cats ADD CONSTRAINT FK_shop_cats_user
		FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_product_id
		FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_filter_value_id
		FOREIGN KEY (filter_value_id) REFERENCES shop_filters_values (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_filters_values ADD CONSTRAINT FK_shop_filters_values_filter_id
		FOREIGN KEY (filter_id) REFERENCES shop_filters (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_user
		FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_currency
		FOREIGN KEY (currency) REFERENCES shop_currencies (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	return nil
}
