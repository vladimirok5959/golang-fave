package modules

import (
	"os"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

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
			if _, err := tx.Exec("SELECT product_id FROM shop_product_images WHERE product_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Delete product attached images
			rows, err := wrap.DB.Query(
				`SELECT
					product_id,
					filename
				FROM
					shop_product_images
				WHERE
					product_id = ?
				;`,
				pf_id,
			)
			if err == nil {
				defer rows.Close()
				values := make([]string, 2)
				scan := make([]interface{}, len(values))
				for i := range values {
					scan[i] = &values[i]
				}
				for rows.Next() {
					err = rows.Scan(scan...)
					if err == nil {
						target_file_full := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + string(values[0]) + string(os.PathSeparator) + string(values[1])
						os.Remove(target_file_full)
					}
				}
			}
			if _, err := tx.Exec("DELETE FROM shop_product_images WHERE product_id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Delete target product with category connection data
			if _, err := tx.Exec("DELETE FROM shop_filter_product_values WHERE product_id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM shop_cat_product_rel WHERE product_id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM shop_products WHERE id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
