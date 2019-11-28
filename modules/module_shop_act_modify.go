package modules

import (
	"errors"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_gname := wrap.R.FormValue("gname")
		pf_name := wrap.R.FormValue("name")
		pf_price := wrap.R.FormValue("price")
		pf_price_old := wrap.R.FormValue("price_old")
		pf_currency := wrap.R.FormValue("currency")
		pf_alias := wrap.R.FormValue("alias")
		pf_vendor := wrap.R.FormValue("vendor")
		pf_quantity := wrap.R.FormValue("quantity")
		pf_category := wrap.R.FormValue("category")
		pf_briefly := wrap.R.FormValue("briefly")
		pf_content := wrap.R.FormValue("content")
		pf_active := wrap.R.FormValue("active")

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsFloat(pf_price) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsFloat(pf_price_old) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_currency) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_quantity) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_category) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify product name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateSingleAlias(pf_name)
		}

		if !utils.IsValidSingleAlias(pf_alias) {
			wrap.MsgError(`Please specify correct product alias`)
			return
		}

		// Default is ROOT
		if pf_category == "0" {
			pf_category = "1"
		}

		// Collect fields and data for filter values
		filter_values := map[int]int{}
		for key, values := range wrap.R.PostForm {
			if len(key) > 6 && key[0:6] == "value." {
				for _, value := range values {
					if value != "" {
						filter_values[utils.StrToInt(value)] = utils.StrToInt(key[6:])
					}
				}
			}
		}

		if pf_id == "0" {
			var lastID int64 = 0
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Insert row
				res, err := tx.Exec(
					`INSERT INTO shop_products SET
						user = ?,
						currency = ?,
						price = ?,
						price_old = ?,
						gname = ?,
						name = ?,
						alias = ?,
						vendor = ?,
						quantity = ?,
						category = ?,
						briefly = ?,
						content = ?,
						datetime = ?,
						active = ?
					;`,
					wrap.User.A_id,
					utils.StrToInt(pf_currency),
					utils.StrToFloat64(pf_price),
					utils.StrToFloat64(pf_price_old),
					pf_gname,
					pf_name,
					pf_alias,
					pf_vendor,
					utils.StrToInt(pf_quantity),
					utils.StrToInt(pf_category),
					pf_briefly,
					pf_content,
					utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
					utils.StrToInt(pf_active),
				)
				if err != nil {
					return err
				}

				// Get inserted product id
				lastID, err = res.LastInsertId()
				if err != nil {
					return err
				}

				// Block rows
				if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", lastID); err != nil {
					return err
				}

				// Insert product and categories relations
				catids := utils.GetPostArrayInt("cats[]", wrap.R)
				if len(catids) > 0 {
					var catsCount int
					err = tx.QueryRow(`
						SELECT
							COUNT(*)
						FROM
							shop_cats
						WHERE
							id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
						FOR UPDATE;`,
					).Scan(
						&catsCount,
					)
					if *wrap.LogCpError(&err) != nil {
						return err
					}
					if len(catids) != catsCount {
						return errors.New("Inner system error")
					}
					var bulkInsertArr []string
					for _, el := range catids {
						bulkInsertArr = append(bulkInsertArr, `(`+utils.Int64ToStr(lastID)+`,`+utils.IntToStr(el)+`)`)
					}
					if _, err = tx.Exec(
						`INSERT INTO shop_cat_product_rel (product_id,category_id) VALUES ` + strings.Join(bulkInsertArr, ",") + `;`,
					); err != nil {
						return err
					}
				}

				// Insert product and filter values relations
				for vid, _ := range filter_values {
					if _, err = tx.Exec(
						`INSERT INTO shop_filter_product_values SET
							product_id = ?,
							filter_value_id = ?
						;`,
						lastID,
						vid,
					); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.RecreateProductXmlFile()

			wrap.ResetCacheBlocks()

			wrap.Write(`window.location='/cp/shop/modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Block rows
				if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}
				if _, err := tx.Exec("SELECT id FROM shop_currencies WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_currency)); err != nil {
					return err
				}
				if _, err := tx.Exec("SELECT product_id FROM shop_cat_product_rel WHERE product_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}
				if _, err := tx.Exec("SELECT product_id FROM shop_filter_product_values WHERE product_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}

				// Update row
				if _, err := tx.Exec(
					`UPDATE shop_products SET
						currency = ?,
						price = ?,
						price_old = ?,
						gname = ?,
						name = ?,
						alias = ?,
						vendor = ?,
						quantity = ?,
						category = ?,
						briefly = ?,
						content = ?,
						active = ?
					WHERE
						id = ?
					;`,
					utils.StrToInt(pf_currency),
					utils.StrToFloat64(pf_price),
					utils.StrToFloat64(pf_price_old),
					pf_gname,
					pf_name,
					pf_alias,
					pf_vendor,
					utils.StrToInt(pf_quantity),
					utils.StrToInt(pf_category),
					pf_briefly,
					pf_content,
					utils.StrToInt(pf_active),
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}

				// Delete product and categories relations
				if _, err := tx.Exec("DELETE FROM shop_cat_product_rel WHERE product_id = ?;", utils.StrToInt(pf_id)); err != nil {
					return err
				}

				// Insert product and categories relations
				catids := utils.GetPostArrayInt("cats[]", wrap.R)
				if len(catids) > 0 {
					var catsCount int
					err := tx.QueryRow(`
						SELECT
							COUNT(*)
						FROM
							shop_cats
						WHERE
							id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
						FOR UPDATE;`,
					).Scan(
						&catsCount,
					)
					if *wrap.LogCpError(&err) != nil {
						return err
					}
					if len(catids) != catsCount {
						return errors.New("Inner system error")
					}
					var bulkInsertArr []string
					for _, el := range catids {
						bulkInsertArr = append(bulkInsertArr, `(`+pf_id+`,`+utils.IntToStr(el)+`)`)
					}
					if _, err := tx.Exec(
						`INSERT INTO shop_cat_product_rel (product_id,category_id) VALUES ` + strings.Join(bulkInsertArr, ",") + `;`,
					); err != nil {
						return err
					}
				}

				// Delete product and filter values relations
				if _, err := tx.Exec(
					`DELETE FROM shop_filter_product_values WHERE product_id = ?;`,
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}

				// Insert product and filter values relations
				for vid, _ := range filter_values {
					if _, err := tx.Exec(
						`INSERT INTO shop_filter_product_values SET
							product_id = ?,
							filter_value_id = ?
						;`,
						utils.StrToInt(pf_id),
						vid,
					); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.RecreateProductXmlFile()

			wrap.ResetCacheBlocks()

			wrap.Write(`window.location='/cp/shop/modify/` + pf_id + `/';`)
		}
	})
}
