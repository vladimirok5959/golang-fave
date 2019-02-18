package modules

import (
	"html/template"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type MInfo struct {
	Id     string
	WantDB bool
	Mount  string
	Name   string
	Order  int
	System bool
}

type Module struct {
	Info  MInfo
	Front func(wrap *wrapper.Wrapper)
	Back  func(wrap *wrapper.Wrapper) (string, string, string)
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

func (this *Modules) newModule(info MInfo, ff func(wrap *wrapper.Wrapper), bf func(wrap *wrapper.Wrapper) (string, string, string)) *Module {
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

func (this *Modules) getNavMenuModules(wrap *wrapper.Wrapper, sys bool) string {
	list := make([]*MInfo, 0)
	for _, mod := range this.mods {
		if mod.Back != nil {
			if mod.Info.System == sys {
				list = append(list, &mod.Info)
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	html := ""
	for _, mod := range list {
		class := ""
		if mod.Mount == wrap.CurrModule {
			class = " active"
		}
		html += `<a class="dropdown-item` + class + `" href="/cp/` + mod.Mount + `/">` + mod.Name + `</a>`
	}
	return html
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
			sidebar_left, content, sidebar_right := mod.Back(wrap)

			body_class := "cp"
			if sidebar_left != "" {
				body_class = body_class + " cp-sidebar-left"
			}
			if content == "" {
				body_class = body_class + " cp-404"
				content = "Panel 404"
			}
			if sidebar_right != "" {
				body_class = body_class + " cp-sidebar-right"
			}

			wrap.RenderBackEnd(assets.TmplCpBase, consts.TmplDataCpBase{
				Title:              "Fave " + consts.ServerVersion,
				BodyClasses:        body_class,
				UserId:             wrap.User.A_id,
				UserFirstName:      wrap.User.A_first_name,
				UserLastName:       wrap.User.A_last_name,
				UserEmail:          wrap.User.A_email,
				UserPassword:       "",
				UserAvatarLink:     "https://s.gravatar.com/avatar/" + utils.GetMd5(wrap.User.A_email) + "?s=80&r=g",
				NavBarModules:      template.HTML(this.getNavMenuModules(wrap, false)),
				NavBarModulesSys:   template.HTML(this.getNavMenuModules(wrap, true)),
				ModuleCurrentAlias: wrap.CurrModule,
				SidebarLeft:        template.HTML(sidebar_left),
				Content:            template.HTML(content),
				SidebarRight:       template.HTML(sidebar_right),
			})

			return true
		}
	}
	return false
}
