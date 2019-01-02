package actions

import (
//"database/sql"
//_ "github.com/go-sql-driver/mysql"

//utils "golang-fave/engine/wrapper/utils"
)

func (this *Action) Action_singout() {
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

	// Set to zero
	this.wrapper.Session.SetInt("UserId", 0)

	// Reload current page
	this.write(`window.location.reload(false);`)
}
