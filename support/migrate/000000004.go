package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000004(db *sqlw.DB) error {
	if _, err := db.Exec(
		`ALTER TABLE blog_cat_post_rel DROP id;`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		`ALTER TABLE shop_cat_product_rel DROP id;`,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		`ALTER TABLE shop_filter_product_values DROP id;`,
	); err != nil {
		return err
	}
	return nil
}
