package actions

import (
	"fmt"
)

func (this *Action) Action_signin() {
	this.msg_success(`Hello from web server`)
}
