package modules

import (
	"fmt"
	"net/http"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Some() *Module {
	return this.newModule(false, "some", "Some Module", func(wrap *wrapper.Wrapper) {
		fmt.Printf("SOME FrontEnd func call\n")
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`SOME FrontEnd func call (` + wrap.CurrModule + `)`))
	}, func(wrap *wrapper.Wrapper) {
		fmt.Printf("SOME BackEnd func call\n")
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`SOME BackEnd func call (` + wrap.CurrModule + `)`))
	})
}
