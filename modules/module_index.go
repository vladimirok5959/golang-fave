package modules

import (
	"fmt"
	"net/http"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Index() *Module {
	return this.newModule(false, "index", "Pages", func(wrap *wrapper.Wrapper) {
		//fmt.Printf("FrontEnd func call\n")
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`FrontEnd func call`))
	}, func(wrap *wrapper.Wrapper) {
		//fmt.Printf("BackEnd func call\n")
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`BackEnd func call`))
	})
}

func (this *Modules) RegisterAction_MysqlSetup() *Action {
	return this.newAction(false, "mysql", func(wrap *wrapper.Wrapper) {
		fmt.Printf("ActFunc func call\n")
	})
}

// All actions here...
// MySQL install
// MySQL first user
// User login
// User logout
