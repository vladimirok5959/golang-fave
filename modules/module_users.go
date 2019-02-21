package modules

import (
	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Users() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "users",
		Name:   "Users",
		Order:  800,
		System: true,
		Icon:   assets.SysSvgIconUser,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of Users", Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add New User", Icon: assets.SysSvgIconPlus},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of Users"},
			})
		} else if wrap.CurrSubModule == "add" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Add New User"},
			})
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
