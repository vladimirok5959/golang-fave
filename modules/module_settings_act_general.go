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
		pf_module_at_home := utils.Trim(wrap.R.FormValue("module-at-home"))

		if !utils.IsNumeric(pf_module_at_home) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_module_at_home := utils.StrToInt(pf_module_at_home)

		// Correct values
		if pfi_module_at_home < 0 {
			pfi_module_at_home = 0
		}
		if pfi_module_at_home > 2 {
			pfi_module_at_home = 2
		}

		(*wrap.Config).Engine.MainModule = pfi_module_at_home

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
