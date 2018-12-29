package actions

import (
	"fmt"
)

func action_mysql(e *Action) {
	/*
	action := e.R.FormValue("action")
	e.write(fmt.Sprintf(`
		ModalShowMsg('MySQL Action', 'Hello from web server (%s)');
	`, action))
	*/

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
	err := e.MySqlConfigWrite(pf_host, pf_name, pf_user, pf_password)
	

	/*
	if pf_host == "" || pf_name == "" || pf_user == "" || pf_password == "" {
		e.write(fmt.Sprintf(`ModalShowMsg('Error', 'Not all fields are filed');`))
		return
	}
	*/

	// Reload current page
	e.write(fmt.Sprintf(`window.location.reload(false);`))
}
