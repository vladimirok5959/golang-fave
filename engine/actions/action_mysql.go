package actions

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	utils "golang-fave/engine/wrapper/utils"
)

func action_mysql(e *Action) {
	pf_host := e.R.FormValue("host")
	pf_name := e.R.FormValue("name")
	pf_user := e.R.FormValue("user")
	pf_password := e.R.FormValue("password")

	if pf_host == "" {
		e.msg_error(`Please specify host for mysql connection`)
		return
	}

	if pf_name == "" {
		e.msg_error(`Please specify mysql database name`)
		return
	}

	if pf_user == "" {
		e.msg_error(`Please specify mysql user`)
		return
	}

	// Try connect to mysql
	db, err := sql.Open("mysql", pf_user+":"+pf_password+"@tcp("+pf_host+":3306)/"+pf_name)
	if err != nil {
		e.msg_error(err.Error())
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		e.msg_error(err.Error())
		return
	}

	// Try to install all tables

	// Save mysql config file
	err = utils.MySqlConfigWrite(e.VHostHome, pf_host, pf_name, pf_user, pf_password)
	if err != nil {
		e.msg_error(err.Error())
		return
	}

	// Reload current page
	e.write(fmt.Sprintf(`window.location.reload(false);`))
}
