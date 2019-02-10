package modules

import (
	//"fmt"
	"net/http"
	"reflect"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type Module struct {
	Id       string
	WantDB   bool
	Mount    string
	Name     string
	FrontEnd func(wrap *wrapper.Wrapper)
	BackEnd  func(wrap *wrapper.Wrapper)
}

type Action struct {
	Id      string
	WantDB  bool
	Mount   string
	ActFunc func(wrap *wrapper.Wrapper)
}

type Modules struct {
	mods map[string]*Module
	acts map[string]*Action
}

func (this *Modules) load() {
	t := reflect.TypeOf(this)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if strings.HasPrefix(m.Name, "XXX") {
			continue
		}
		if strings.HasPrefix(m.Name, "RegisterModule_") {
			id := m.Name[15:]
			if _, ok := reflect.TypeOf(this).MethodByName("RegisterModule_" + id); ok {
				result := reflect.ValueOf(this).MethodByName("RegisterModule_" + id).Call([]reflect.Value{})
				if len(result) >= 1 {
					mod := result[0].Interface().(*Module)
					mod.Id = id
					this.mods[mod.Mount] = mod
				}
			}
		}
		if strings.HasPrefix(m.Name, "RegisterAction_") {
			id := m.Name[15:]
			if _, ok := reflect.TypeOf(this).MethodByName("RegisterAction_" + id); ok {
				result := reflect.ValueOf(this).MethodByName("RegisterAction_" + id).Call([]reflect.Value{})
				if len(result) >= 1 {
					act := result[0].Interface().(*Action)
					act.Id = id
					this.acts[act.Mount] = act
				}
			}
		}
	}
}

func (this *Modules) newModule(WantDB bool, Mount string, Name string, ff func(wrap *wrapper.Wrapper), bf func(wrap *wrapper.Wrapper)) *Module {
	return &Module{
		WantDB:   WantDB,
		Mount:    Mount,
		Name:     Name,
		FrontEnd: ff,
		BackEnd:  bf,
	}
}

func (this *Modules) newAction(WantDB bool, Mount string, af func(wrap *wrapper.Wrapper)) *Action {
	return &Action{
		WantDB:  WantDB,
		Mount:   Mount,
		ActFunc: af,
	}
}

func New() *Modules {
	m := Modules{
		mods: map[string]*Module{},
		acts: map[string]*Action{},
	}
	m.load()
	return &m
}

func (this *Modules) XXXActionFire(wrap *wrapper.Wrapper) bool {
	if wrap.R.Method == "POST" {
		if err := wrap.R.ParseForm(); err == nil {
			name := wrap.R.FormValue("action")
			if name != "" {
				wrap.W.WriteHeader(http.StatusOK)
				wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				wrap.W.Header().Set("Content-Type", "text/html; charset=utf-8")
				if act, ok := this.acts[name]; ok {
					if act.WantDB {
						err := wrap.UseDatabase()
						if err != nil {
							wrap.MsgError(err.Error())
							return true
						}
						defer wrap.DB.Close()
					}
					act.ActFunc(wrap)
					return true
				} else {
					wrap.MsgError(`This action is not implemented`)
					return true
				}
			}
		}
	}
	return false
}

func (this *Modules) XXXFrontEnd(wrap *wrapper.Wrapper) bool {
	var mod *Module = nil

	// Some module
	if len(wrap.UrlArgs) > 0 {
		if m, ok := this.mods[wrap.UrlArgs[0]]; ok {
			mod = m
		}
	}

	// Default module
	if mod == nil {
		if m, ok := this.mods["index"]; ok {
			mod = m
		}
	}

	// Check and run
	if mod != nil {
		if mod.FrontEnd != nil {
			if mod.WantDB {
				err := wrap.UseDatabase()
				if err != nil {
					utils.SystemErrorPageEngine(wrap.W, err)
					return true
				}
				defer wrap.DB.Close()
			}
			mod.FrontEnd(wrap)
			return true
		}
	}

	return false
}

func (this *Modules) XXXBackEnd(wrap *wrapper.Wrapper) bool {
	//fmt.Printf("Back: %v\n", wrap.UrlArgs)
	return false
}
