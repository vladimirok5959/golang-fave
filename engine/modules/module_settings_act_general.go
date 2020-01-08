package modules

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsGeneral() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-general",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_module_at_home := utils.Trim(wrap.R.FormValue("module-at-home"))
		pf_maintenance := utils.Trim(wrap.R.FormValue("maintenance"))

		if !utils.IsNumeric(pf_module_at_home) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if pf_maintenance == "" {
			pf_maintenance = "0"
		}
		if !utils.IsNumeric(pf_maintenance) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_module_at_home := utils.StrToInt(pf_module_at_home)
		pfi_maintenance := utils.StrToInt(pf_maintenance)

		// Correct values
		if pfi_module_at_home < 0 {
			pfi_module_at_home = 0
		}
		if pfi_module_at_home > 2 {
			pfi_module_at_home = 2
		}

		if pfi_maintenance < 0 {
			pfi_maintenance = 0
		}
		if pfi_maintenance > 1 {
			pfi_maintenance = 1
		}

		(*wrap.Config).Engine.MainModule = pfi_module_at_home
		(*wrap.Config).Engine.Maintenance = pfi_maintenance

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
