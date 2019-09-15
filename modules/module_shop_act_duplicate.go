package modules

import (
	"time"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopDuplicate() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-duplicate",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		var lastID int64 = 0
		if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT product_id FROM shop_cat_product_rel WHERE product_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT product_id FROM shop_filter_product_values WHERE product_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Duplicate product
			res, err := tx.Exec(
				`INSERT INTO shop_products (
					user,
					currency,
					price,
					name,
					alias,
					vendor,
					quantity,
					category,
					briefly,
					content,
					datetime,
					active
				) SELECT
					user,
					currency,
					price,
					CONCAT(name, ' (Copy)'),
					CONCAT(alias, '-', 'copy-`+utils.Int64ToStr(time.Now().Unix())+`'),
					vendor,
					quantity,
					category,
					briefly,
					content,
					datetime,
					0
				FROM
					shop_products
				WHERE
					id = ?
				;`,
				utils.StrToInt(pf_id),
			)
			if err != nil {
				return err
			}

			// Get inserted product id
			lastID, err = res.LastInsertId()
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Delete products XML cache
		wrap.RemoveProductXmlCacheFile()

		wrap.ResetCacheBlocks()

		// Navigate to new product page
		wrap.Write(`window.location='/cp/shop/modify/` + utils.Int64ToStr(lastID) + `/';`)
	})
}
