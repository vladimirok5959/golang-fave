package modules

import (
	others "golang-fave/engine/wrapper/resources/others"
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

func (this *Module) Module_index_icon() string {
	return others.File_assets_sys_svg_page
}

func (this *Module) Module_index_order() int {
	return 1
}

func (this *Module) Module_index_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{
		Alias: "default",
		Name:  "List of pages",
		Icon:  others.File_assets_sys_svg_list,
	})
	result = append(result, utils.ModuleSubMenu{
		Alias: "modify",
		Name:  "Add new page",
		Icon:  others.File_assets_sys_svg_plus,
	})
	return result
}

func (this *Module) Module_index_content() string {
	return "Index content"
}

func (this *Module) Module_index_sidebar() string {
	return "Index right sidebar"
}
