package modules

import (
	"context"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_ShopCategoriesDelete() *Action {
	return this.newAction(AInfo{
		Mount:     "shop-categories-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))

		if !utils.IsNumeric(pf_id) || utils.StrToInt(pf_id) <= 1 {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_shop_cats FOR UPDATE;"); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT category_id FROM fave_shop_cat_product_rel WHERE category_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_shop_products WHERE category = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Set root category
			if _, err := tx.Exec(ctx, "UPDATE fave_shop_products SET category = 1 WHERE category = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec(ctx, "DELETE FROM fave_shop_cat_product_rel WHERE category_id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT @ml := lft, @mr := rgt FROM fave_shop_cats WHERE id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "DELETE FROM fave_shop_cats WHERE id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "UPDATE fave_shop_cats SET lft = lft - 1, rgt = rgt - 1 WHERE lft > @ml AND rgt < @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "UPDATE fave_shop_cats SET lft = lft - 2 WHERE lft > @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "UPDATE fave_shop_cats SET rgt = rgt - 2 WHERE rgt > @mr;"); err != nil {
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
