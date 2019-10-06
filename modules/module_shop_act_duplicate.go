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
		pf_attach := wrap.R.FormValue("attach")

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

			parent_id := "parent_id"
			if pf_attach == "1" {
				parent_id = pf_id
			}

			// Duplicate product
			res, err := tx.Exec(
				`INSERT INTO shop_products (
					parent_id,
					user,
					currency,
					price,
					gname,
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
					`+parent_id+`,
					user,
					currency,
					price,
					'',
					CONCAT(name, ' (Copy)'),
					CONCAT(REGEXP_REPLACE(alias, '-c[0-9]+$', ''), '-c', '`+utils.Int64ToStr(time.Now().Unix())+`'),
					vendor,
					quantity,
					category,
					briefly,
					content,
					'`+utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp())+`',
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

			// Block new product row
			if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", lastID); err != nil {
				return err
			}

			// Duplicate categories
			cat_sqls := []string{}
			if cat_rows, err := tx.Query(
				`SELECT
					product_id,
					category_id
				FROM
					shop_cat_product_rel
				WHERE
					product_id = ?
				;`,
				utils.StrToInt(pf_id),
			); err == nil {
				defer cat_rows.Close()
				for cat_rows.Next() {
					var product_id int
					var category_id int
					if err := cat_rows.Scan(&product_id, &category_id); *wrap.LogCpError(&err) == nil {
						cat_sqls = append(cat_sqls, `
							INSERT INTO shop_cat_product_rel SET
								product_id = `+utils.Int64ToStr(lastID)+`,
								category_id = `+utils.IntToStr(category_id)+`
							;
						`)
					}
				}
			}
			for _, sql_query := range cat_sqls {
				tx.Exec(sql_query)
			}

			// Duplicate attributes
			attributes_sqls := []string{}
			if attributes_rows, err := tx.Query(
				`SELECT
					product_id,
					filter_value_id
				FROM
					shop_filter_product_values
				WHERE
					product_id = ?
				;`,
				utils.StrToInt(pf_id),
			); err == nil {
				defer attributes_rows.Close()
				for attributes_rows.Next() {
					var product_id int
					var filter_value_id int
					if err := attributes_rows.Scan(&product_id, &filter_value_id); *wrap.LogCpError(&err) == nil {
						attributes_sqls = append(attributes_sqls, `
							INSERT INTO shop_filter_product_values SET
								product_id = `+utils.Int64ToStr(lastID)+`,
								filter_value_id = `+utils.IntToStr(filter_value_id)+`
							;
						`)
					}
				}
			}
			for _, sql_query := range attributes_sqls {
				tx.Exec(sql_query)
			}

			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Navigate to new product page
		wrap.Write(`window.location='/cp/shop/modify/` + utils.Int64ToStr(lastID) + `/';`)
	})
}
