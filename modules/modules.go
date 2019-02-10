package modules

import (
	"fmt"

	"golang-fave/engine/wrapper"
)

type FuncFrontEnd func(wrap *wrapper.Wrapper)
type FuncBackEnd func(wrap *wrapper.Wrapper)

type Module struct {
	Id       string
	Name     string
	FrontEnd FuncFrontEnd
	BackEnd  FuncBackEnd
}

type Modules struct {
	list map[string]Module
}

func New() *Modules {
	m := Modules{
		list: map[string]Module{},
	}
	return &m
}

func (this *Modules) Load() {
	// Called before server starts
	fmt.Println("Load modules")
}

// All actions here...
// MySQL install
// MySQL first user
// User login
// User logout

// Called inside goroutine
func (this *Modules) Action(wrap *wrapper.Wrapper) bool {
	//
	return false
}

// Called inside goroutine
func (this *Modules) FrontEnd(wrap *wrapper.Wrapper) bool {
	//
	return false
}

// Called inside goroutine
func (this *Modules) BackEnd(wrap *wrapper.Wrapper) bool {
	//
	return false
}
