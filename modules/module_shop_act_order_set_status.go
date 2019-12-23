package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopOrderSetStatus() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-order-set-status",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_status := utils.Trim(wrap.R.FormValue("status"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_status) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !(utils.StrToInt(pf_status) >= 0 && utils.StrToInt(pf_status) <= 4) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if err := wrap.DB.Transaction(wrap.R.Context(), func(tx *wrapper.Tx) error {
			if _, err := tx.Exec(
				`UPDATE shop_orders SET
					status = ?
				WHERE
					id = ?
				;`,
				utils.StrToInt(pf_status),
				utils.StrToInt(pf_id),
			); err != nil {
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
