package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopCurrenciesDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-currencies-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) || utils.StrToInt(pf_id) <= 1 {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Process
			if _, err := tx.Exec(
				`DELETE FROM shop_currencies WHERE id = ?;`,
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
