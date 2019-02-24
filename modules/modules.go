package modules

import (
	"html"
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

type MISub struct {
	Mount string
	Name  string
	Icon  string
	Show  bool
}

type MInfo struct {
	Id     string
	WantDB bool
	Mount  string
	Name   string
	Order  int
	System bool
	Icon   string
	Sub    *[]MISub
}

type Module struct {
	Info  MInfo
	Front func(wrap *wrapper.Wrapper)
	Back  func(wrap *wrapper.Wrapper) (string, string, string)
}

type AInfo struct {
	Id        string
	WantDB    bool
	Mount     string
	WantUser  bool
	WantAdmin bool
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
	if !backend || (backend && len(wrap.UrlArgs) <= 0) {
		if mod == nil {
			if m, ok := this.mods["index"]; ok {
				mod = m
				modCurr = "index"
			}
		}
	}

	return mod, modCurr
}

func (this *Modules) getModulesList(wrap *wrapper.Wrapper, sys bool, all bool) []*MInfo {
	list := make([]*MInfo, 0)
	for _, mod := range this.mods {
		if mod.Back != nil {
			if mod.Info.System == sys || all {
				list = append(list, &mod.Info)
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	return list
}

func (this *Modules) getSidebarModuleSubMenu(wrap *wrapper.Wrapper, mod *MInfo) string {
	html := ``
	if mod.Sub != nil {
		for _, item := range *mod.Sub {
			if item.Show {
				class := ""
				if (item.Mount == "default" && len(wrap.UrlArgs) <= 1) || (len(wrap.UrlArgs) >= 2 && item.Mount == wrap.UrlArgs[1]) {
					class = " active"
				}
				icon := item.Icon
				if icon == "" {
					icon = assets.SysSvgIconGear
				}
				href := "/cp/" + mod.Mount + "/" + item.Mount + "/"
				if mod.Mount == "index" && item.Mount == "default" {
					href = "/cp/"
				} else if item.Mount == "default" {
					href = "/cp/" + mod.Mount + "/"
				}
				html += `<li class="nav-item` + class + `"><a class="nav-link" href="` + href + `">` + icon + item.Name + `</a></li>`
			}
		}
		if html != "" {
			html = `<ul class="nav flex-column">` + html + `</ul>`
		}
	}
	return html
}

func (this *Modules) getNavMenuModules(wrap *wrapper.Wrapper, sys bool) string {
	html := ``
	list := this.getModulesList(wrap, sys, false)
	for _, mod := range list {
		class := ""
		if mod.Mount == wrap.CurrModule {
			class = " active"
		}
		href := `/cp/` + mod.Mount + `/`
		if mod.Mount == "index" {
			href = `/cp/`
		}
		html += `<a class="dropdown-item` + class + `" href="` + href + `">` + mod.Name + `</a>`
	}
	return html
}

func (this *Modules) getSidebarModules(wrap *wrapper.Wrapper) string {
	html_def := ""
	html_sys := ""
	list := this.getModulesList(wrap, false, true)
	for _, mod := range list {
		class := ""
		submenu := ""
		if mod.Mount == wrap.CurrModule {
			class = " active"
			submenu = this.getSidebarModuleSubMenu(wrap, mod)
		}
		icon := mod.Icon
		if icon == "" {
			icon = assets.SysSvgIconGear
		}
		href := "/cp/" + mod.Mount + "/"
		if mod.Mount == "index" {
			href = "/cp/"
		}
		if !mod.System {
			html_def += `<li class="nav-item` + class + `"><a class="nav-link" href="` + href + `">` + icon + mod.Name + `</a>` + submenu + `</li>`
		} else {
			html_sys += `<li class="nav-item` + class + `"><a class="nav-link" href="` + href + `">` + icon + mod.Name + `</a>` + submenu + `</li>`
		}
	}
	if html_def != "" {
		html_def = `<ul class="nav flex-column">` + html_def + `</ul>`
	}
	if html_sys != "" {
		html_sys = `<ul class="nav flex-column">` + html_sys + `</ul>`
	}
	if html_def != "" && html_sys != "" {
		html_sys = `<div class="dropdown-divider"></div>` + html_sys
	}
	return html_def + html_sys
}

func (this *Modules) getBreadCrumbs(wrap *wrapper.Wrapper, data *[]consts.BreadCrumb) string {
	res := `<nav aria-label="breadcrumb">`
	res += `<ol class="breadcrumb">`
	res += `<li class="breadcrumb-item"><a href="/cp/` + this.mods[wrap.CurrModule].Info.Mount + `/">` + html.EscapeString(this.mods[wrap.CurrModule].Info.Name) + `</a></li>`
	for _, item := range *data {
		if item.Link == "" {
			res += `<li class="breadcrumb-item active" aria-current="page">` + html.EscapeString(item.Name) + `</li>`
		} else {
			res += `<li class="breadcrumb-item"><a href="` + item.Link + `">` + html.EscapeString(item.Name) + `</a></li>`
		}
	}
	res += `</ol>`
	res += `</nav>`
	return res
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
					if act.Info.WantUser || act.Info.WantAdmin {
						if !wrap.LoadSessionUser() {
							wrap.MsgError(`You must be loginned to run this action`)
							return true
						}
						if wrap.User.A_active <= 0 {
							if !wrap.LoadSessionUser() {
								wrap.MsgError(`You do not have rights to run this action`)
								return true
							}
						}
					}
					if act.Info.WantAdmin && wrap.User.A_admin <= 0 {
						if !wrap.LoadSessionUser() {
							wrap.MsgError(`You do not have rights to run this action`)
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
		if len(wrap.UrlArgs) >= 2 && wrap.UrlArgs[1] != "" {
			wrap.CurrSubModule = wrap.UrlArgs[1]
		}

		// Search for sub module mount
		found := false
		submount := "default"
		if wrap.CurrSubModule != "" {
			submount = wrap.CurrSubModule
		}
		for _, item := range *mod.Info.Sub {
			if item.Mount == submount {
				found = true
				break
			}
		}

		// Display standart 404 error page
		if !found {
			return found
		}

		// Call module function
		if mod.Back != nil {
			sidebar_left, content, sidebar_right := mod.Back(wrap)

			// Display standart 404 error page
			if sidebar_left == "" && content == "" && sidebar_right == "" {
				return false
			}

			// Prepare CP page
			body_class := "cp"
			if sidebar_left != "" {
				body_class = body_class + " cp-sidebar-left"
			}
			if content == "" {
				body_class = body_class + " cp-404"
			}
			if sidebar_right != "" {
				body_class = body_class + " cp-sidebar-right"
			}

			wrap.RenderBackEnd(assets.TmplCpBase, consts.TmplDataCpBase{
				Title:              wrap.CurrHost + " - Fave " + consts.ServerVersion,
				Caption:            "Fave " + consts.ServerVersion,
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
