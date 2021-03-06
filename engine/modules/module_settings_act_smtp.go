package modules

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsSmtp() *Action {
	return this.newAction(AInfo{
		Mount:     "settings-smtp",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_smtp_host := utils.Trim(wrap.R.FormValue("smtp-host"))
		pf_smtp_port := utils.Trim(wrap.R.FormValue("smtp-port"))
		pf_smtp_login := utils.Trim(wrap.R.FormValue("smtp-login"))
		pf_smtp_password := utils.Trim(wrap.R.FormValue("smtp-password"))
		pf_smtp_test_email := utils.Trim(wrap.R.FormValue("smtp-test-email"))

		(*wrap.Config).SMTP.Host = pf_smtp_host
		(*wrap.Config).SMTP.Port = utils.StrToInt(pf_smtp_port)
		(*wrap.Config).SMTP.Login = pf_smtp_login

		// Update password only if not empty
		if pf_smtp_password != "" {
			(*wrap.Config).SMTP.Password = pf_smtp_password
		}

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Send test message
		if pf_smtp_test_email != "" {
			if err := wrap.SendEmailFast(
				pf_smtp_test_email,
				"❤️ Fave.Pro SMTP test message",
				"Hello! This is Fave.Pro test message.<br />If you see this message, then you right configured SMTP settings!",
			); err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
