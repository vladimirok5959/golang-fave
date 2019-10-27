package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_SettingsGeneral() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-general",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_module_at_home := wrap.R.FormValue("module-at-home")
		pf_price_fomat := wrap.R.FormValue("price-fomat")
		pf_price_round := wrap.R.FormValue("price-round")

		if !utils.IsNumeric(pf_module_at_home) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if !utils.IsNumeric(pf_price_fomat) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if !utils.IsNumeric(pf_price_round) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_module_at_home := utils.StrToInt(pf_module_at_home)
		pfi_price_fomat := utils.StrToInt(pf_price_fomat)
		pfi_price_round := utils.StrToInt(pf_price_round)

		// Correct values
		if pfi_module_at_home < 0 {
			pfi_module_at_home = 0
		}
		if pfi_module_at_home > 2 {
			pfi_module_at_home = 2
		}

		if pfi_price_fomat < 0 {
			pfi_price_fomat = 0
		}
		if pfi_price_fomat > 4 {
			pfi_price_fomat = 4
		}

		if pfi_price_round < 0 {
			pfi_price_round = 0
		}
		if pfi_price_round > 2 {
			pfi_price_round = 2
		}

		(*wrap.Config).Engine.MainModule = pfi_module_at_home
		(*wrap.Config).Shop.Price.Format = pfi_price_fomat
		(*wrap.Config).Shop.Price.Round = pfi_price_round

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
