package modules

import (
	"golang-fave/engine/wrapper"
)

type Modules struct {
	//
}

func New() *Modules {
	m := Modules{}
	return &m
}

func (this *Modules) Load() {
	// Called before server starts
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
