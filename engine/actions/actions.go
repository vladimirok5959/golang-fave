package actions

import (
	"net/http"
)

type hRun func(e *Action)

type Action struct {
	W         *http.ResponseWriter
	R         *http.Request
	VHost     string
	VHostHome string
	RemoteIp  string
	list      map[string]hRun
}

func (e *Action) register(name string, handle hRun) {
	e.list[name] = handle
}

func (e *Action) write(data string) {
	(*e.W).Write([]byte(data))
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, vhosthome string, remoteip string) *Action {
	act := Action{w, r, vhost, vhosthome, remoteip, make(map[string]hRun)}

	// Register all action here
	act.register("mysql", action_mysql)
	act.register("signin", action_signin)

	return &act
}

func (e *Action) Call() bool {
	if e.R.Method != "POST" {
		return false
	}
	if err := e.R.ParseForm(); err == nil {
		action := e.R.FormValue("action")
		if action != "" {
			fn, ok := e.list[action]
			if ok {
				(*e.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				(*e.W).Header().Set("Content-Type", "text/html; charset=utf-8")
				fn(e)
				return true
			}
		}
	}
	return false
}
