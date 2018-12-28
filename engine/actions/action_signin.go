package actions

import (
	"fmt"
)

func action_signin(e *Action) {
	action := e.R.FormValue("action")
	(*e.W).Write([]byte(fmt.Sprintf(`
		alert('Action: %s');
	`, action)))
}
