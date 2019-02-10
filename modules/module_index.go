package modules

import (
	"fmt"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_index() *Module {
	return this.newModule(false, "Pages", func(wrap *wrapper.Wrapper) {
		fmt.Printf("FrontEnd func call\n")
	}, func(wrap *wrapper.Wrapper) {
		fmt.Printf("BackEnd func call\n")
	})
}

func (this *Modules) RegisterAction_mysql() *Action {
	return this.newAction(false, func(wrap *wrapper.Wrapper) {
		fmt.Printf("ActFunc func call\n")
	})
}

// All actions here...
// MySQL install
// MySQL first user
// User login
// User logout
