package modules

import (
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Module) Module_index() {
	// Do something here...
}

func (this *Module) Module_index_display() bool {
	return true
}

func (this *Module) Module_index_alias() string {
	return "index"
}

func (this *Module) Module_index_name() string {
	return "Pages"
}

func (this *Module) Module_index_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{Alias: "default", Name: "List of pages"})
	result = append(result, utils.ModuleSubMenu{Alias: "add", Name: "Add new page"})
	return result
}

func (this *Module) Module_index_content() string {
	return "Index content"
}

func (this *Module) Module_index_sidebar() string {
	return "Index right sidebar"
}
