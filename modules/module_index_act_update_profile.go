package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_IndexUserUpdateProfile() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "index-user-update-profile",
		WantUser: true,
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

		if pf_password != "" {
			// Update with password if set
			_, err := wrap.DB.Query(
				`UPDATE users SET
					first_name = ?,
					last_name = ?,
					email = ?,
					password = MD5(?)
				WHERE
					id = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				pf_password,
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		} else {
			// Update without password if not set
			_, err := wrap.DB.Query(
				`UPDATE users SET
					first_name = ?,
					last_name = ?,
					email = ?
				WHERE
					id = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
