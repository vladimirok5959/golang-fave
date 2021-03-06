package modules

import (
	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_IndexUserSignIn() *Action {
	return this.newAction(AInfo{
		Mount: "index-user-sign-in",
	}, func(wrap *wrapper.Wrapper) {
		pf_email := utils.Trim(wrap.R.FormValue("email"))
		pf_password := utils.Trim(wrap.R.FormValue("password"))

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

		if wrap.S.GetInt("UserId", 0) > 0 {
			wrap.MsgError(`You already logined`)
			return
		}

		var user_id int
		err := wrap.DB.QueryRow(
			wrap.R.Context(),
			`SELECT
				id
			FROM
				fave_users
			WHERE
				email = ? and
				password = MD5(?) and
				admin = 1 and
				active = 1
			LIMIT 1;`,
			pf_email,
			pf_password,
		).Scan(
			&user_id,
		)

		if err != nil && err != sqlw.ErrNoRows {
			wrap.LogCpError(&err)
			wrap.MsgError(err.Error())
			return
		}

		if err == sqlw.ErrNoRows {
			wrap.MsgError(`Incorrect email or password`)
			return
		}

		// Save to current session
		wrap.S.SetInt("UserId", user_id)

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
