package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_UsersModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "users-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")
		pf_admin := wrap.R.FormValue("admin")
		pf_active := wrap.R.FormValue("active")

		if pf_admin == "" {
			pf_admin = "0"
		}

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		// First user always super admin
		// Rewrite active and admin status
		if pf_id == "1" {
			pf_admin = "1"
			pf_active = "1"
		}

		if pf_id == "0" {
			// Add new user
			if pf_password == "" {
				wrap.MsgError(`Please specify user password`)
				return
			}
			_, err := wrap.DB.Exec(
				`INSERT INTO users SET
					first_name = ?,
					last_name = ?,
					email = ?,
					password = MD5(?),
					admin = ?,
					active = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				pf_password,
				pf_admin,
				pf_active,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/users/';`)
		} else {
			// Update user
			if pf_password == "" {
				_, err := wrap.DB.Exec(
					`UPDATE users SET
						first_name = ?,
						last_name = ?,
						email = ?,
						admin = ?,
						active = ?
					WHERE
						id = ?
					;`,
					pf_first_name,
					pf_last_name,
					pf_email,
					pf_admin,
					pf_active,
					utils.StrToInt(pf_id),
				)
				if err != nil {
					wrap.MsgError(err.Error())
					return
				}
			} else {
				_, err := wrap.DB.Exec(
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
					utils.StrToInt(pf_id),
				)
				if err != nil {
					wrap.MsgError(err.Error())
					return
				}
			}
			wrap.Write(`window.location='/cp/users/modify/` + pf_id + `/';`)
		}
	})
}
