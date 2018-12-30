package actions

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	utils "golang-fave/engine/wrapper/utils"
)

func (this *Action) Action_mysql() {
	pf_host := this.wrapper.R.FormValue("host")
	pf_port := this.wrapper.R.FormValue("port")
	pf_name := this.wrapper.R.FormValue("name")
	pf_user := this.wrapper.R.FormValue("user")
	pf_password := this.wrapper.R.FormValue("password")

	if pf_host == "" {
		this.msg_error(`Please specify host for MySQL connection`)
		return
	}

	if pf_port == "" {
		this.msg_error(`Please specify host port for MySQL connection`)
		return
	}

	if _, err := strconv.Atoi(pf_port); err != nil {
		this.msg_error(`MySQL host port must be integer number`)
		return
	}

	if pf_name == "" {
		this.msg_error(`Please specify MySQL database name`)
		return
	}

	if pf_user == "" {
		this.msg_error(`Please specify MySQL user`)
		return
	}

	// Try connect to mysql
	db, err := sql.Open("mysql", pf_user+":"+pf_password+"@tcp("+pf_host+":"+pf_port+")/"+pf_name)
	if err != nil {
		this.msg_error(err.Error())
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Try to install all tables

	// Save mysql config file
	err = utils.MySqlConfigWrite(this.wrapper.DirVhostHome, pf_host, pf_port, pf_name, pf_user, pf_password)
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Reload current page
	this.write(fmt.Sprintf(`window.location.reload(false);`))
}
