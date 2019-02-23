package modules

import (
	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Settings() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "settings",
		Name:   "Settings",
		Order:  801,
		System: true,
		Icon:   assets.SysSvgIconGear,
		Sub: &[]MISub{
			{Mount: "default", Name: "Main Settings", Show: true, Icon: assets.SysSvgIconList},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Main Settings"},
			})
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
