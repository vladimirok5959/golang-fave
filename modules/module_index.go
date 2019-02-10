package modules

import (
	"fmt"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_index() *Module {
	return &Module{
		Id:   "index",
		Name: "Pages",
		FrontEnd: func(mod *Modules, wrap *wrapper.Wrapper) {
			fmt.Printf("FrontEnd func call\n")
		},
		BackEnd: func(mod *Modules, wrap *wrapper.Wrapper) {
			fmt.Printf("BackEnd func call\n")
		},
	}
}

func (this *Modules) RegisterAction_test() *Action {
	return &Action{
		Id: "test",
		ActFunc: func(mod *Modules, wrap *wrapper.Wrapper) {
			fmt.Printf("ActFunc func call\n")
		},
	}
}

/*
func (this *Modules) module_index_frontend(wrap *wrapper.Wrapper) {
	//
}

func (this *Modules) module_index_backend(wrap *wrapper.Wrapper) {
	//
}
*/
