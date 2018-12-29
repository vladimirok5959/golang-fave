package actions

import (
	"fmt"
)

func action_mysql(e *Action) {
	action := e.R.FormValue("action")
	(*e.W).Write([]byte(fmt.Sprintf(`
		ModalShowMsg('MySQL Action', 'Hello from web server (%s)');
	`, action)))
}
