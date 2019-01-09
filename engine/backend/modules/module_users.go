package modules

import (
	others "golang-fave/engine/wrapper/resources/others"
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Module) Module_users() {
	// Do something here...
}

func (this *Module) Module_users_display() bool {
	return false
}

func (this *Module) Module_users_alias() string {
	return "users"
}

func (this *Module) Module_users_name() string {
	return "Users"
}

func (this *Module) Module_users_icon() string {
	return others.File_assets_sys_svg_user
}

func (this *Module) Module_users_order() int {
	return 100
}

func (this *Module) Module_users_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{
		Alias: "default",
		Name:  "List of users",
		Icon:  others.File_assets_sys_svg_list,
	})
	result = append(result, utils.ModuleSubMenu{
		Alias: "modify",
		Name:  "Add new user",
		Icon:  others.File_assets_sys_svg_plus,
	})
	return result
}

func (this *Module) Module_users_content() string {
	return "Users content"
}

func (this *Module) Module_users_sidebar() string {
	return "Users right sidebar"
}
