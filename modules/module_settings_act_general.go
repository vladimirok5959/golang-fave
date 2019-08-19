package modules

import (
	"strconv"

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

		if _, err := strconv.Atoi(pf_module_at_home); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_module_at_home := utils.StrToInt(pf_module_at_home)

		// Correct some values
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

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
