package modules

import (
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
	return "Empty"
}

func (this *Module) Module_empty_order() int {
	return 999
}

func (this *Module) Module_empty_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{Alias: "default", Name: "Some list"})
	return result
}

func (this *Module) Module_empty_content() string {
	return "Empty content"
}

func (this *Module) Module_empty_sidebar() string {
	return "Empty right sidebar"
}
