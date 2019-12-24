package modules

import (
	"context"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopDetach() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-detach",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			if _, err := tx.Exec(
				ctx,
				`UPDATE shop_products SET
					parent_id = NULL,
					active = 0
				WHERE
					id = ?
				;`,
				utils.StrToInt(pf_id),
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
