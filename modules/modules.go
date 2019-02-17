package modules

import (
	"net/http"
	"reflect"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type MInfo struct {
	Id     string
	WantDB bool
	Mount  string
	Name   string
}

type Module struct {
	Info  MInfo
	Front func(wrap *wrapper.Wrapper)
	Back  func(wrap *wrapper.Wrapper)
}

type AInfo struct {
	Id       string
	WantDB   bool
	Mount    string
	WantUser bool
}

type Action struct {
	Info AInfo
	Act  func(wrap *wrapper.Wrapper)
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
					mod.Info.Id = id
					this.mods[mod.Info.Mount] = mod
				}
			}
		}
		if strings.HasPrefix(m.Name, "RegisterAction_") {
			id := m.Name[15:]
			if _, ok := reflect.TypeOf(this).MethodByName("RegisterAction_" + id); ok {
				result := reflect.ValueOf(this).MethodByName("RegisterAction_" + id).Call([]reflect.Value{})
				if len(result) >= 1 {
					act := result[0].Interface().(*Action)
					act.Info.Id = id
					this.acts[act.Info.Mount] = act
				}
			}
		}
	}
}

func (this *Modules) newModule(info MInfo, ff func(wrap *wrapper.Wrapper), bf func(wrap *wrapper.Wrapper)) *Module {
	return &Module{Info: info, Front: ff, Back: bf}
}

func (this *Modules) newAction(info AInfo, af func(wrap *wrapper.Wrapper)) *Action {
	return &Action{Info: info, Act: af}
}

func (this *Modules) getCurrentModule(wrap *wrapper.Wrapper, backend bool) (*Module, string) {
	var mod *Module = nil
	var modCurr string = ""

	// Some module
	if len(wrap.UrlArgs) >= 1 {
		if m, ok := this.mods[wrap.UrlArgs[0]]; ok {
			if (!backend && m.Front != nil) || (backend && m.Back != nil) {
				mod = m
				modCurr = wrap.UrlArgs[0]
			}
		}
	}

	// Default module
	if mod == nil {
		if m, ok := this.mods["index"]; ok {
			mod = m
			modCurr = "index"
		}
	}

	return mod, modCurr
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
					if act.Info.WantDB {
						err := wrap.UseDatabase()
						if err != nil {
							wrap.MsgError(err.Error())
							return true
						}
						defer wrap.DB.Close()
					}
					if act.Info.WantUser {
						if !wrap.LoadSessionUser() {
							wrap.MsgError(`You must be loginned to run this action`)
							return true
						}
					}
					act.Act(wrap)
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
	mod, cm := this.getCurrentModule(wrap, false)
	if mod != nil {
		wrap.CurrModule = cm
		if mod.Front != nil {
			if mod.Info.WantDB {
				err := wrap.UseDatabase()
				if err != nil {
					utils.SystemErrorPageEngine(wrap.W, err)
					return true
				}
				defer wrap.DB.Close()
			}
			mod.Front(wrap)
			return true
		}
	}
	return false
}

func (this *Modules) XXXBackEnd(wrap *wrapper.Wrapper) bool {
	mod, cm := this.getCurrentModule(wrap, true)
	if mod != nil {
		wrap.CurrModule = cm
		if mod.Back != nil {
			mod.Back(wrap)
			return true
		}
	}
	return false
}
