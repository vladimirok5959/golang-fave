package modules

import (
	"context"
	"errors"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopAttachProductTo() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-attach-product-to",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_parent_id := utils.Trim(wrap.R.FormValue("parent_id"))
		pf_product_id := utils.Trim(wrap.R.FormValue("product_id"))

		if !utils.IsNumeric(pf_parent_id) || !utils.IsNumeric(pf_product_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			// Check parent
			var count int
			if err := tx.QueryRow(
				ctx,
				"SELECT COUNT(*) FROM `fave_shop_products` WHERE `id` = ? AND `parent_id` IS NULL;",
				utils.StrToInt(pf_parent_id),
			).Scan(&count); err != nil {
				return err
			}
			if count <= 0 {
				return errors.New("Parent product can't be used for attaching")
			}

			// Check child
			if err := tx.QueryRow(
				ctx,
				"SELECT COUNT(*) FROM `fave_shop_products` WHERE `parent_id` = ?;",
				utils.StrToInt(pf_product_id),
			).Scan(&count); err != nil {
				return err
			}
			if count >= 1 {
				return errors.New("Parent can't be attached to parent")
			}

			// Attach
			if _, err := tx.Exec(
				ctx,
				`UPDATE fave_shop_products SET
					parent_id = ?
				WHERE
					id = ? AND
					parent_id IS NULL
				;`,
				utils.StrToInt(pf_parent_id),
				utils.StrToInt(pf_product_id),
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.RecreateProductXmlFile()

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
