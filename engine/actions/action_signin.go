package actions

import (
	"fmt"
)

func action_signin(e *Action) {
	action := e.w.R.FormValue("action")
	e.write(fmt.Sprintf(`
		ModalShowMsg('Login Action', 'Hello from web server (%s)');
	`, action))
}
