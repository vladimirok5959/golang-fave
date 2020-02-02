package modules

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsGeneral() *Action {
	return this.newAction(AInfo{
		Mount:     "settings-general",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_module_at_home := utils.Trim(wrap.R.FormValue("module-at-home"))
		pf_maintenance := utils.Trim(wrap.R.FormValue("maintenance"))
		pf_mod_enabled_blog := utils.Trim(wrap.R.FormValue("mod-enabled-blog"))
		pf_mod_enabled_shop := utils.Trim(wrap.R.FormValue("mod-enabled-shop"))

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

		if pf_mod_enabled_blog == "" {
			pf_mod_enabled_blog = "0"
		}
		if !utils.IsNumeric(pf_mod_enabled_blog) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if pf_mod_enabled_shop == "" {
			pf_mod_enabled_shop = "0"
		}
		if !utils.IsNumeric(pf_mod_enabled_shop) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_module_at_home := utils.StrToInt(pf_module_at_home)
		pfi_maintenance := utils.StrToInt(pf_maintenance)
		pfi_mod_enabled_blog := utils.StrToInt(pf_mod_enabled_blog)
		pfi_mod_enabled_shop := utils.StrToInt(pf_mod_enabled_shop)

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

		if pfi_mod_enabled_blog < 0 {
			pfi_mod_enabled_blog = 0
		}
		if pfi_mod_enabled_blog > 1 {
			pfi_mod_enabled_blog = 1
		}

		if pfi_mod_enabled_shop < 0 {
			pfi_mod_enabled_shop = 0
		}
		if pfi_mod_enabled_shop > 1 {
			pfi_mod_enabled_shop = 1
		}

		(*wrap.Config).Engine.MainModule = pfi_module_at_home
		(*wrap.Config).Engine.Maintenance = pfi_maintenance
		(*wrap.Config).Modules.Enabled.Blog = pfi_mod_enabled_blog
		(*wrap.Config).Modules.Enabled.Shop = pfi_mod_enabled_shop

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
