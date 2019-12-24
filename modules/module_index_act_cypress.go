package modules

import (
	"os"

	"golang-fave/consts"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_IndexCypressReset() *Action {
	return this.newAction(AInfo{
		WantDB: false,
		Mount:  "index-cypress-reset",
	}, func(wrap *wrapper.Wrapper) {
		if !consts.ParamDebug {
			wrap.Write(`Access denied`)
			return
		}

		db, err := sqlw.Open("mysql", "root:root@tcp(localhost:3306)/fave")
		if err != nil {
			wrap.Write(err.Error())
			return
		}
		defer db.Close()
		err = db.Ping(wrap.R.Context())
		if err != nil {
			wrap.Write(err.Error())
			return
		}

		os.Remove(wrap.DConfig + string(os.PathSeparator) + "mysql.json")
		os.Remove(wrap.DConfig + string(os.PathSeparator) + "config.json")
		wrap.RemoveProductImageThumbnails("*", "*")

		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE blog_cat_post_rel DROP FOREIGN KEY FK_blog_cat_post_rel_post_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE blog_cat_post_rel DROP FOREIGN KEY FK_blog_cat_post_rel_category_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE blog_cats DROP FOREIGN KEY FK_blog_cats_user;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE blog_posts DROP FOREIGN KEY FK_blog_posts_user;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE blog_posts DROP FOREIGN KEY FK_blog_posts_category;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE pages DROP FOREIGN KEY FK_pages_user;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_cat_product_rel DROP FOREIGN KEY FK_shop_cat_product_rel_product_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_cat_product_rel DROP FOREIGN KEY FK_shop_cat_product_rel_category_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_cats DROP FOREIGN KEY FK_shop_cats_user;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_filter_product_values DROP FOREIGN KEY FK_shop_filter_product_values_product_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_filter_product_values DROP FOREIGN KEY FK_shop_filter_product_values_filter_value_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_filters_values DROP FOREIGN KEY FK_shop_filters_values_filter_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_orders DROP FOREIGN KEY FK_shop_orders_currency_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_order_products DROP FOREIGN KEY FK_shop_order_products_order_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_order_products DROP FOREIGN KEY FK_shop_order_products_product_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_product_images DROP FOREIGN KEY FK_shop_product_images_product_id;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_user;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_currency;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_category;`)
		_, _ = db.Exec(wrap.R.Context(), `ALTER TABLE shop_products DROP FOREIGN KEY FK_shop_products_parent_id;`)

		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE blog_cat_post_rel;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE blog_cats;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE blog_posts;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE notify_mail;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE pages;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE settings;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_cat_product_rel;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_cats;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_currencies;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_filter_product_values;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_filters_values;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_filters;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_order_products;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_orders;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_product_images;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE shop_products;`)
		_, _ = db.Exec(wrap.R.Context(), `DROP TABLE users;`)

		wrap.ResetCacheBlocks()

		wrap.Write(`OK`)
	})
}
