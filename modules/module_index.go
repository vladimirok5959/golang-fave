package modules

import (
	"net/http"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Index() *Module {
	return this.newModule(false, "index", "Pages", func(wrap *wrapper.Wrapper) {
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`INDEX FrontEnd func call (` + wrap.CurrModule + `)`))
	}, func(wrap *wrapper.Wrapper) {
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`INDEX BackEnd func call (` + wrap.CurrModule + `)`))
	})
}

func (this *Modules) RegisterAction_MysqlSetup() *Action {
	return this.newAction(false, "mysql", func(wrap *wrapper.Wrapper) {
		wrap.MsgError(`Some error`)
	})
}

// All actions here...
// MySQL install
// MySQL first user
// User login
// User logout
