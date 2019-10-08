package modules

import (
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_SettingsSmtp() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-smtp",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_smtp_host := strings.TrimSpace(wrap.R.FormValue("smtp-host"))
		pf_smtp_port := strings.TrimSpace(wrap.R.FormValue("smtp-port"))
		pf_smtp_login := strings.TrimSpace(wrap.R.FormValue("smtp-login"))
		pf_smtp_password := strings.TrimSpace(wrap.R.FormValue("smtp-password"))

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

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
