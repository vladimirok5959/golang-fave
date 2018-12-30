package actions

import (
	"fmt"
	"regexp"
)

var regexpe = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (this *Action) Action_first_user() {
	if dbe := this.use_database(); dbe != nil {
		this.msg_error(dbe.Error())
		return
	} else {
		defer this.db.Close()
	}

	pf_first_name := this.wrapper.R.FormValue("first_name")
	pf_last_name := this.wrapper.R.FormValue("last_name")
	pf_email := this.wrapper.R.FormValue("email")
	pf_password := this.wrapper.R.FormValue("password")

	if pf_email == "" {
		this.msg_error(`Please specify user email`)
		return
	}

	if !regexpe.MatchString(pf_email) {
		this.msg_error(`Please specify correct user email`)
		return
	}

	if pf_password == "" {
		this.msg_error(`Please specify user password`)
		return
	}

	_, err := this.db.Query(
		"INSERT INTO `users` SET `first_name` = ?, `last_name` = ?, `email` = ?, `password` = MD5(?);",
		pf_first_name, pf_last_name, pf_email, pf_password)
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Reload current page
	this.write(fmt.Sprintf(`window.location.reload(false);`))
}
