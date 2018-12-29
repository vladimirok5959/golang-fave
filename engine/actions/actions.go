package actions

import (
	"fmt"
	"strings"

	"golang-fave/engine/wrapper"
)

type hRun func(e *Action)

type Action struct {
	w    *wrapper.Wrapper
	list map[string]hRun
}

func (e *Action) register(name string, handle hRun) {
	e.list[name] = handle
}

func (e *Action) write(data string) {
	(*e.w.W).Write([]byte(data))
}

func (e *Action) msg_show(title string, msg string) {
	e.write(fmt.Sprintf(
		`ModalShowMsg('%s', '%s');`,
		strings.Replace(strings.Replace(title, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1),
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}

func (e *Action) msg_success(msg string) {
	e.msg_show("Success", msg)
}

func (e *Action) msg_error(msg string) {
	e.msg_show("Error", msg)
}

func New(w *wrapper.Wrapper) *Action {
	act := Action{w, make(map[string]hRun)}

	// Register all action here
	act.register("mysql", action_mysql)
	act.register("signin", action_signin)

	return &act
}

func (e *Action) Call() bool {
	if e.w.R.Method != "POST" {
		return false
	}
	if err := e.w.R.ParseForm(); err == nil {
		action := e.w.R.FormValue("action")
		if action != "" {
			function, ok := e.list[action]
			if ok {
				(*e.w.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				(*e.w.W).Header().Set("Content-Type", "text/html; charset=utf-8")
				function(e)
				return true
			}
		}
	}
	return false
}
