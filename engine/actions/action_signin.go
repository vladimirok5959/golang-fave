package actions

func (this *Action) Action_signin() {
	if dbe := this.use_database(); dbe != nil {
		this.msg_error(dbe.Error())
		return
	} else {
		defer this.db.Close()
	}

	this.msg_success(`Hello from web server`)
}
