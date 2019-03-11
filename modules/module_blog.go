package modules

import (
	//"html"

	"golang-fave/assets"
	"golang-fave/consts"
	//"golang-fave/engine/builder"
	"golang-fave/engine/wrapper"
	//"golang-fave/utils"
)

func (this *Modules) RegisterModule_Blog() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "blog",
		Name:   "Blog",
		Order:  1,
		System: false,
		Icon:   assets.SysSvgIconPage,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of posts", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new post", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify post", Show: false},
			{Mount: "cats", Name: "List of categories", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "cats-add", Name: "Add new category", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "cats-modify", Name: "Modify category", Show: false},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of posts"},
			})
			//
		} else if wrap.CurrSubModule == "cats" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of categories"},
			})
			//
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new post"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify post"},
				})
			}
			//
		} else if wrap.CurrSubModule == "cats-add" || wrap.CurrSubModule == "cats-modify" {
			if wrap.CurrSubModule == "cats-add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new category"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify category"},
				})
			}
			//
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
