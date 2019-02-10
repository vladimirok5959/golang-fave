package modules

import (
	"fmt"
	"reflect"
	"strings"

	"golang-fave/engine/wrapper"
)

type FuncFrontEnd func(mod *Modules, wrap *wrapper.Wrapper)
type FuncBackEnd func(mod *Modules, wrap *wrapper.Wrapper)
type FuncAction func(mod *Modules, wrap *wrapper.Wrapper)

type Module struct {
	Id       string
	Name     string
	FrontEnd FuncFrontEnd
	BackEnd  FuncBackEnd
}

type Action struct {
	Id      string
	ActFunc FuncAction
}

type Modules struct {
	mods map[string]Module
	acts map[string]Action
}

func New() *Modules {
	m := Modules{
		mods: map[string]Module{},
		acts: map[string]Action{},
	}
	m.load()
	return &m
}

func (this *Modules) load() {
	// Called before server starts
	fmt.Println("Load modules")
	fmt.Println("---")

	t := reflect.TypeOf(this)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if strings.HasPrefix(m.Name, "XXX") {
			continue
		}
		if strings.HasPrefix(m.Name, "RegisterModule_") {
			//fmt.Printf("%s\n", m.Name)
			id := m.Name[15:]
			fmt.Printf("Module (%s)\n", id)

			if _, ok := reflect.TypeOf(this).MethodByName("RegisterModule_" + id); ok {
				result := reflect.ValueOf(this).MethodByName("RegisterModule_" + id).Call([]reflect.Value{})
				if len(result) >= 1 {
					//fmt.Printf("--- result -> (%d)\n", len(result))
					//fmt.Printf("%v\n", result[0])
					//fmt.Printf("%T\n", result)
					//mod := result[0]
					//mod := result[0]
					//fmt.Printf("%v\n", mod)
					//fmt.Printf("%s\n", *Module(mod).Id)
					mod := result[0].Interface().(*Module)
					fmt.Printf("%s\n", mod.Id)
				}
			} else {
				fmt.Printf("Error\n")
			}
			// Add to array
			//mod := 
			/*
			if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
				result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
				return result[0].String()
			}
			*/
		}
		if strings.HasPrefix(m.Name, "RegisterAction_") {
			//fmt.Printf("%s\n", m.Name)
			id := m.Name[15:]
			fmt.Printf("Action (%s)\n", id)
			// Add to array
		}
	}
}

// All actions here...
// MySQL install
// MySQL first user
// User login
// User logout

// Called inside goroutine
func (this *Modules) XXXActionFire(wrap *wrapper.Wrapper) bool {
	//
	return false
}

// Called inside goroutine
func (this *Modules) XXXFrontEnd(wrap *wrapper.Wrapper) bool {
	//
	return false
}

// Called inside goroutine
func (this *Modules) XXXBackEnd(wrap *wrapper.Wrapper) bool {
	//
	return false
}
