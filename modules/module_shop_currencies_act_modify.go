package modules

import (
	"context"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopCurrenciesModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-currencies-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_name := utils.Trim(wrap.R.FormValue("name"))
		pf_coefficient := utils.Trim(wrap.R.FormValue("coefficient"))
		pf_code := utils.Trim(wrap.R.FormValue("code"))
		pf_symbol := utils.Trim(wrap.R.FormValue("symbol"))

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
			var lastID int64 = 0
			if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
				// Insert row
				res, err := tx.Exec(
					ctx,
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

				// Get inserted id
				lastID, err = res.LastInsertId()
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.RecreateProductXmlFile()

			wrap.ResetCacheBlocks()

			wrap.Write(`window.location='/cp/shop/currencies-modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
				// Block rows
				if _, err := tx.Exec(ctx, "SELECT id FROM shop_currencies WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}

				// Update row
				if _, err := tx.Exec(
					ctx,
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

			wrap.RecreateProductXmlFile()

			wrap.ResetCacheBlocks()

			wrap.Write(`window.location='/cp/shop/currencies-modify/` + pf_id + `/';`)
		}
	})
}
