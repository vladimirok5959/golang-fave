package actions

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

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

	var user_id int
	err := this.db.QueryRow(
		"SELECT `id` FROM `users` WHERE `email` = ? and `password` = MD5(?) LIMIT 1;",
		pf_email, pf_password).Scan(&user_id)

	if err != nil && err != sql.ErrNoRows {
		this.msg_error(err.Error())
		return
	}

	if err == sql.ErrNoRows {
		this.msg_error(`Incorrect email or password`)
		return
	}

	/*
	if !this.Session.IsSetInt("UserId") {
		this.Session.SetInt("UserId", 0)
	}
	if !this.Session.IsSetBool("IsLogged") {
		this.Session.SetBool("IsLogged", false)
	}
	*/

	var session_user_id int
	session_user_id, _ = this.wrapper.Session.GetInt("UserId")

	this.msg_success(fmt.Sprintf(
		`Test: (%d), (%d)`,
		user_id, session_user_id))

	// Reload current page
	//this.write(fmt.Sprintf(`window.location.reload(false);`))

	//this.msg_success(`Hello from web server`)
	//this.write(fmt.Sprintf(``))

	/*
	if count <= 0 {
		return this.wrapper.TmplBackEnd(templates.CpFirstUser, nil)
	}
	*/

	//this.msg_success(`Hello from web server`)
}
