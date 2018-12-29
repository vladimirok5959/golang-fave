package actions

import (
	"fmt"

	utils "golang-fave/engine/wrapper/utils"
)

func action_mysql(e *Action) {
	pf_host := e.R.FormValue("host")
	pf_name := e.R.FormValue("name")
	pf_user := e.R.FormValue("user")
	pf_password := e.R.FormValue("password")

	if pf_host == "" {
		e.write(fmt.Sprintf(`ModalShowMsg('Error', 'Please specify host for mysql connection');`))
		return
	}

	if pf_name == "" {
		e.write(fmt.Sprintf(`ModalShowMsg('Error', 'Please specify mysql database name');`))
		return
	}

	if pf_user == "" {
		e.write(fmt.Sprintf(`ModalShowMsg('Error', 'Please specify mysql user');`))
		return
	}

	// Try connect to mysql

	// Try to install all tables

	// Save mysql config file
	err := utils.MySqlConfigWrite(e.VHostHome, pf_host, pf_name, pf_user, pf_password)
	if err != nil {
		e.write(fmt.Sprintf(`ModalShowMsg('Error', '%s');`, err))
	}

	// Reload current page
	e.write(fmt.Sprintf(`window.location.reload(false);`))
}
