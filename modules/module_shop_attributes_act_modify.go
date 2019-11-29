package modules

import (
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopAttributesModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-attributes-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_name := utils.Trim(wrap.R.FormValue("name"))
		pf_filter := utils.Trim(wrap.R.FormValue("filter"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify attribute name`)
			return
		}

		if pf_filter == "" {
			wrap.MsgError(`Please specify attribute in filter`)
			return
		}

		// Collect fields and data
		filter_values := map[string]int{}
		for key, values := range wrap.R.PostForm {
			if len(key) > 6 && key[0:6] == "value." {
				for _, value := range values {
					if value != "" {
						filter_values[value] = utils.StrToInt(key[6:])
					}
				}
			}
		}

		if pf_id == "0" {
			var lastID int64 = 0
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Insert row
				res, err := tx.Exec(
					`INSERT INTO shop_filters SET
						name = ?,
						filter = ?
					;`,
					pf_name,
					pf_filter,
				)
				if err != nil {
					return err
				}

				// Get inserted id
				lastID, err = res.LastInsertId()
				if err != nil {
					return err
				}

				// Block rows
				if _, err := tx.Exec("SELECT id FROM shop_filters WHERE id = ? FOR UPDATE;", lastID); err != nil {
					return err
				}

				// Insert values
				for vname, _ := range filter_values {
					if _, err = tx.Exec(
						`INSERT INTO shop_filters_values SET
							filter_id = ?,
							name = ?
						;`,
						lastID,
						vname,
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

			wrap.Write(`window.location='/cp/shop/attributes-modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Block rows
				if _, err := tx.Exec("SELECT id FROM shop_filters WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}
				if _, err := tx.Exec("SELECT id FROM shop_filters_values WHERE filter_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}
				if _, err := tx.Exec(
					`SELECT
						shop_filter_product_values.product_id
					FROM
						shop_filter_product_values
						LEFT JOIN shop_filters_values ON shop_filters_values.id = shop_filter_product_values.filter_value_id
					WHERE
						shop_filters_values.id IS NOT NULL AND
						shop_filters_values.filter_id = ?
					FOR UPDATE;`,
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}

				// Update row
				if _, err := tx.Exec(
					`UPDATE shop_filters SET
						name = ?,
						filter = ?
					WHERE
						id = ?
					;`,
					pf_name,
					pf_filter,
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}

				// Delete not existed rows
				ignore_ids := []string{}
				for _, vid := range filter_values {
					if vid != 0 {
						ignore_ids = append(ignore_ids, utils.IntToStr(vid))
					}
				}
				if len(ignore_ids) > 0 {
					if _, err := tx.Exec(
						`DELETE
							shop_filter_product_values
						FROM
							shop_filter_product_values
							LEFT JOIN shop_filters_values ON shop_filters_values.id = shop_filter_product_values.filter_value_id
						WHERE
							shop_filters_values.id IS NOT NULL AND
							shop_filters_values.filter_id = ? AND
							shop_filter_product_values.filter_value_id NOT IN (`+strings.Join(ignore_ids, ",")+`)
						;`,
						utils.StrToInt(pf_id),
					); err != nil {
						return err
					}
					if _, err := tx.Exec(
						`DELETE FROM shop_filters_values WHERE filter_id = ? AND id NOT IN (`+strings.Join(ignore_ids, ",")+`);`,
						utils.StrToInt(pf_id),
					); err != nil {
						return err
					}
				} else {
					if _, err := tx.Exec(
						`DELETE
							shop_filter_product_values
						FROM
							shop_filter_product_values
							LEFT JOIN shop_filters_values ON shop_filters_values.id = shop_filter_product_values.filter_value_id
						WHERE
							shop_filters_values.id IS NOT NULL AND
							shop_filters_values.filter_id = ?
						;`,
						utils.StrToInt(pf_id),
					); err != nil {
						return err
					}
					if _, err := tx.Exec(
						`DELETE FROM shop_filters_values WHERE filter_id = ?;`,
						utils.StrToInt(pf_id),
					); err != nil {
						return err
					}
				}

				// Insert new values, update existed rows
				for vname, vid := range filter_values {
					if vid == 0 {
						if _, err := tx.Exec(
							`INSERT INTO shop_filters_values SET
								filter_id = ?,
								name = ?
							;`,
							utils.StrToInt(pf_id),
							vname,
						); err != nil {
							return err
						}
					} else {
						if _, err := tx.Exec(
							`UPDATE shop_filters_values SET
								name = ?
							WHERE
								id = ? AND
								filter_id = ?
							;`,
							vname,
							vid,
							utils.StrToInt(pf_id),
						); err != nil {
							return err
						}
					}
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.RecreateProductXmlFile()

			wrap.ResetCacheBlocks()

			wrap.Write(`window.location='/cp/shop/attributes-modify/` + pf_id + `/';`)
		}
	})
}
