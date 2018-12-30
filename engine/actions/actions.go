package actions

import (
	"fmt"
	"reflect"
	"strings"

	"golang-fave/engine/wrapper"
)

type Action struct {
	wrapper *wrapper.Wrapper
}

func (this *Action) write(data string) {
	(*this.wrapper.W).Write([]byte(data))
}

func (this *Action) msg_show(title string, msg string) {
	this.write(fmt.Sprintf(
		`ModalShowMsg('%s', '%s');`,
		strings.Replace(strings.Replace(title, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1),
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}

func (this *Action) msg_success(msg string) {
	this.msg_show("Success", msg)
}

func (this *Action) msg_error(msg string) {
	this.msg_show("Error", msg)
}

func New(wrapper *wrapper.Wrapper) *Action {
	return &Action{wrapper}
}

func (this *Action) Run() bool {
	if this.wrapper.R.Method != "POST" {
		return false
	}
	if err := this.wrapper.R.ParseForm(); err == nil {
		action := this.wrapper.R.FormValue("action")
		if action != "" {
			if _, ok := reflect.TypeOf(this).MethodByName("Action_" + action); ok {
				(*this.wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				(*this.wrapper.W).Header().Set("Content-Type", "text/html; charset=utf-8")
				reflect.ValueOf(this).MethodByName("Action_" + action).Call([]reflect.Value{})
				return true
			}
		}
	}
	return false
}
