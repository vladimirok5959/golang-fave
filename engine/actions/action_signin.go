package actions

import (
	"fmt"
)

func (this *Action) Action_signin() {
	action := this.wrapper.R.FormValue("action")
	this.write(fmt.Sprintf(`
		ModalShowMsg('Login Action', 'Hello from web server (%s)');
	`, action))
}
