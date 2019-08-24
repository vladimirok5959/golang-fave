package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopAttributesDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-attributes-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
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

			// Process
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
			if _, err := tx.Exec(
				`DELETE FROM shop_filters WHERE id = ?;`,
				utils.StrToInt(pf_id),
			); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Delete products XML cache
		wrap.RemoveProductXmlCacheFile()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
