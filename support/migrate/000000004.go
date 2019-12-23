package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000004(ctx context.Context, db *sqlw.DB, host string) error {
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE blog_cat_post_rel DROP id;`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_cat_product_rel DROP id;`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		ctx,
		`ALTER TABLE shop_filter_product_values DROP id;`,
	); err != nil {
		return err
	}
	return nil
}
