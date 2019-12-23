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

		_, _ = db.Exec(
			`DROP TABLE
				blog_cat_post_rel,
				blog_cats,
				blog_posts,
				notify_mail,
				pages,
				settings,
				shop_cat_product_rel,
				shop_cats,
				shop_currencies,
				shop_filter_product_values,
				shop_filters,
				shop_filters_values,
				shop_order_products,
				shop_orders,
				shop_product_images,
				shop_products,
				users
			;`,
		)

		wrap.ResetCacheBlocks()

		wrap.Write(`OK`)
	})
}
