package modules

import (
	others "golang-fave/engine/wrapper/resources/others"
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Module) Module_empty() {
	// Do something here...
}

func (this *Module) Module_empty_display() bool {
	return true
}

func (this *Module) Module_empty_alias() string {
	return "empty"
}

func (this *Module) Module_empty_name() string {
	return "Empty module"
}

func (this *Module) Module_empty_icon() string {
	return others.File_assets_sys_svg_gear
}

func (this *Module) Module_empty_order() int {
	return 999
}

func (this *Module) Module_empty_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{
		Alias: "default",
		Name:  "Sub-menu item",
		Icon:  others.File_assets_sys_svg_list,
	})
	return result
}

func (this *Module) Module_empty_content() string {
	return "Empty module content"
}

func (this *Module) Module_empty_sidebar() string {
	return "Empty module right sidebar"
}
