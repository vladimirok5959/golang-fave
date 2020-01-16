package modules

import (
	"context"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_ShopAttributesDelete() *Action {
	return this.newAction(AInfo{
		Mount:     "shop-attributes-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_shop_filters WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_shop_filters_values WHERE filter_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(
				ctx,
				`SELECT
					fave_shop_filter_product_values.product_id
				FROM
					fave_shop_filter_product_values
					LEFT JOIN fave_shop_filters_values ON fave_shop_filters_values.id = fave_shop_filter_product_values.filter_value_id
				WHERE
					fave_shop_filters_values.id IS NOT NULL AND
					fave_shop_filters_values.filter_id = ?
				FOR UPDATE;`,
				utils.StrToInt(pf_id),
			); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec(
				ctx,
				`DELETE
					fave_shop_filter_product_values
				FROM
					fave_shop_filter_product_values
					LEFT JOIN fave_shop_filters_values ON fave_shop_filters_values.id = fave_shop_filter_product_values.filter_value_id
				WHERE
					fave_shop_filters_values.id IS NOT NULL AND
					fave_shop_filters_values.filter_id = ?
				;`,
				utils.StrToInt(pf_id),
			); err != nil {
				return err
			}
			if _, err := tx.Exec(
				ctx,
				`DELETE FROM fave_shop_filters_values WHERE filter_id = ?;`,
				utils.StrToInt(pf_id),
			); err != nil {
				return err
			}
			if _, err := tx.Exec(
				ctx,
				`DELETE FROM fave_shop_filters WHERE id = ?;`,
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

		wrap.RecreateProductXmlFile()

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
