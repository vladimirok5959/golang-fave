package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_IndexFirstUser() *Action {
	return this.newAction(AInfo{
		WantDB: true,
		Mount:  "index-first-user",
	}, func(wrap *wrapper.Wrapper) {
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		if pf_password == "" {
			wrap.MsgError(`Please specify user password`)
			return
		}

		// Security, check if still need to run this action
		var count int
		err := wrap.DB.QueryRow(`
			SELECT
				COUNT(*)
			FROM
				users
			;`,
		).Scan(
			&count,
		)
		if wrap.LogCpError(err) != nil {
			wrap.MsgError(err.Error())
			return
		}
		if count > 0 {
			wrap.MsgError(`CMS is already configured`)
			return
		}

		_, err = wrap.DB.Exec(
			`INSERT INTO users SET
				id = 1,
				first_name = ?,
				last_name = ?,
				email = ?,
				password = MD5(?),
				admin = 1,
				active = 1
			;`,
			pf_first_name,
			pf_last_name,
			pf_email,
			pf_password,
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
