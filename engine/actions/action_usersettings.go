package actions

import (
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Action) Action_usersettings() {
	if err := this.use_database(); err != nil {
		this.msg_error(err.Error())
		return
	} else {
		defer this.db.Close()
	}

	if err := this.load_session_user(); err != nil {
		this.msg_error(err.Error())
		return
	}

	pf_first_name := this.wrapper.R.FormValue("first_name")
	pf_last_name := this.wrapper.R.FormValue("last_name")
	pf_email := this.wrapper.R.FormValue("email")
	pf_password := this.wrapper.R.FormValue("password")

	if pf_email == "" {
		this.msg_error(`Please specify user email`)
		return
	}

	if !utils.EmailIsValid(pf_email) {
		this.msg_error(`Please specify correct user email`)
		return
	}

	if pf_password != "" {
		// Update with password if set
		_, err := this.db.Query(
			"UPDATE `users` SET `first_name` = ?, `last_name` = ?, `email` = ?, `password` = MD5(?) WHERE `id` = ?;",
			pf_first_name, pf_last_name, pf_email, pf_password, this.user.A_id)
		if err != nil {
			this.msg_error(err.Error())
			return
		}
	} else {
		// Update without password if not set
		_, err := this.db.Query(
			"UPDATE `users` SET `first_name` = ?, `last_name` = ?, `email` = ? WHERE `id` = ?;",
			pf_first_name, pf_last_name, pf_email, this.user.A_id)
		if err != nil {
			this.msg_error(err.Error())
			return
		}
	}

	// Reload current page
	this.write(`window.location.reload(false);`)
}
