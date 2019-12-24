package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000021(ctx context.Context, db *sqlw.DB, host string) error {
	// Drop foreign keys
	if _, err := db.Exec(ctx, `ALTER TABLE blog_cat_post_rel DROP FOREIGN KEY FK_blog_cat_post_rel_post_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE blog_cat_post_rel DROP FOREIGN KEY FK_blog_cat_post_rel_category_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE blog_cats DROP FOREIGN KEY FK_blog_cats_user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE blog_posts DROP FOREIGN KEY FK_blog_posts_user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE blog_posts DROP FOREIGN KEY FK_blog_posts_category;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE pages DROP FOREIGN KEY FK_pages_user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_cat_product_rel DROP FOREIGN KEY FK_shop_cat_product_rel_product_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_cat_product_rel DROP FOREIGN KEY FK_shop_cat_product_rel_category_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_cats DROP FOREIGN KEY FK_shop_cats_user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_filter_product_values DROP FOREIGN KEY FK_shop_filter_product_values_product_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_filter_product_values DROP FOREIGN KEY FK_shop_filter_product_values_filter_value_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_filters_values DROP FOREIGN KEY FK_shop_filters_values_filter_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_orders DROP FOREIGN KEY FK_shop_orders_currency_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_order_products DROP FOREIGN KEY FK_shop_order_products_order_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_order_products DROP FOREIGN KEY FK_shop_order_products_product_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_product_images DROP FOREIGN KEY FK_shop_product_images_product_id;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_user;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_currency;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_category;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_parent_id;`); err != nil {
		return err
	}

	// Rename tables
	if _, err := db.Exec(ctx, `RENAME TABLE blog_cat_post_rel TO fave_blog_cat_post_rel;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE blog_cats TO fave_blog_cats;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE blog_posts TO fave_blog_posts;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE notify_mail TO fave_notify_mail;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE pages TO fave_pages;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_cat_product_rel TO fave_shop_cat_product_rel;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_cats TO fave_shop_cats;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_currencies TO fave_shop_currencies;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_filter_product_values TO fave_shop_filter_product_values;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_filters_values TO fave_shop_filters_values;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_filters TO fave_shop_filters;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_order_products TO fave_shop_order_products;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_orders TO fave_shop_orders;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_product_images TO fave_shop_product_images;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE shop_products TO fave_shop_products;`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `RENAME TABLE users TO fave_users;`); err != nil {
		return err
	}

	// Add constraints
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_post_id
		FOREIGN KEY (post_id) REFERENCES fave_blog_posts (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_category_id
		FOREIGN KEY (category_id) REFERENCES fave_blog_cats (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_blog_cats ADD CONSTRAINT FK_blog_cats_user
		FOREIGN KEY (user) REFERENCES fave_users (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_blog_posts ADD CONSTRAINT FK_blog_posts_user
		FOREIGN KEY (user) REFERENCES fave_users (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_blog_posts ADD CONSTRAINT FK_blog_posts_category
		FOREIGN KEY (category) REFERENCES fave_blog_cats (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_pages ADD CONSTRAINT FK_pages_user
		FOREIGN KEY (user) REFERENCES fave_users (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_product_id
		FOREIGN KEY (product_id) REFERENCES fave_shop_products (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_category_id
		FOREIGN KEY (category_id) REFERENCES fave_shop_cats (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_cats ADD CONSTRAINT FK_shop_cats_user
		FOREIGN KEY (user) REFERENCES fave_users (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_product_id
		FOREIGN KEY (product_id) REFERENCES fave_shop_products (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_filter_value_id
		FOREIGN KEY (filter_value_id) REFERENCES fave_shop_filters_values (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_filters_values ADD CONSTRAINT FK_shop_filters_values_filter_id
		FOREIGN KEY (filter_id) REFERENCES fave_shop_filters (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_orders ADD CONSTRAINT FK_shop_orders_currency_id
		FOREIGN KEY (currency_id) REFERENCES fave_shop_currencies (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_order_products ADD CONSTRAINT FK_shop_order_products_order_id
		FOREIGN KEY (order_id) REFERENCES fave_shop_orders (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_order_products ADD CONSTRAINT FK_shop_order_products_product_id
		FOREIGN KEY (product_id) REFERENCES fave_shop_products (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_product_images ADD CONSTRAINT FK_shop_product_images_product_id
		FOREIGN KEY (product_id) REFERENCES fave_shop_products (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_products ADD CONSTRAINT FK_shop_products_user
		FOREIGN KEY (user) REFERENCES fave_users (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_products ADD CONSTRAINT FK_shop_products_currency
		FOREIGN KEY (currency) REFERENCES fave_shop_currencies (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_products ADD CONSTRAINT FK_shop_products_category
		FOREIGN KEY (category) REFERENCES fave_shop_cats (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		ctx,
		`ALTER TABLE fave_shop_products ADD CONSTRAINT FK_shop_products_parent_id
		FOREIGN KEY (parent_id) REFERENCES fave_shop_products (id) ON DELETE RESTRICT;`,
	); err != nil {
		return err
	}

	return nil
}
