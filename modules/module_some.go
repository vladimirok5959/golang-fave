package modules

import (
	"fmt"
	"net/http"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Some() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "some",
		Name:   "Some Module",
	}, nil, func(wrap *wrapper.Wrapper) {
		// Back-end
		fmt.Printf("SOME BackEnd func call\n")
		wrap.W.WriteHeader(http.StatusOK)
		wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrap.W.Write([]byte(`SOME BackEnd func call (` + wrap.CurrModule + `)`))
	})
}
