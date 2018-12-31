package actions

func (this *Action) Action_signin() {
	if dbe := this.use_database(); dbe != nil {
		this.msg_error(dbe.Error())
		return
	} else {
		defer this.db.Close()
	}

	pf_email := this.wrapper.R.FormValue("email")
	pf_password := this.wrapper.R.FormValue("password")

	if pf_email == "" {
		this.msg_error(`Please specify user email`)
		return
	}

	if !this.is_valid_email(pf_email) {
		this.msg_error(`Please specify correct user email`)
		return
	}

	if pf_password == "" {
		this.msg_error(`Please specify user password`)
		return
	}

	//this.msg_success(`Hello from web server`)
}
