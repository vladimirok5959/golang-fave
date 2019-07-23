package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopCurrenciesModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-currencies-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_coefficient := wrap.R.FormValue("coefficient")
		pf_code := wrap.R.FormValue("code")
		pf_symbol := wrap.R.FormValue("symbol")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify currency name`)
			return
		}

		if !utils.IsFloat(pf_coefficient) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_code == "" {
			wrap.MsgError(`Please specify currency code`)
			return
		}

		if pf_symbol == "" {
			wrap.MsgError(`Please specify currency symbol`)
			return
		}

		if pf_id == "0" {
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Insert row
				_, err := tx.Exec(
					`INSERT INTO shop_currencies SET
						name = ?,
						coefficient = ?,
						code = ?,
						symbol = ?
					;`,
					pf_name,
					pf_coefficient,
					pf_code,
					pf_symbol,
				)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.Write(`window.location='/cp/shop/currencies/';`)
		} else {
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Block rows
				if _, err := tx.Exec("SELECT id FROM shop_currencies WHERE id = ? FOR UPDATE;", pf_id); err != nil {
					return err
				}

				// Update row
				if _, err := tx.Exec(
					`UPDATE shop_currencies SET
						name = ?,
						coefficient = ?,
						code = ?,
						symbol = ?
					WHERE
						id = ?
					;`,
					pf_name,
					pf_coefficient,
					pf_code,
					pf_symbol,
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.Write(`window.location='/cp/shop/currencies-modify/` + pf_id + `/';`)
		}
	})
}
