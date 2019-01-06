package modules

import (
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Module) Module_users() {
	// Do something here...
}

func (this *Module) Module_users_alias() string {
	return "users"
}

func (this *Module) Module_users_name() string {
	return "Users"
}

func (this *Module) Module_users_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{Alias: "default", Name: "List of users"})
	result = append(result, utils.ModuleSubMenu{Alias: "add", Name: "Add new user"})
	return result
}

func (this *Module) Module_users_content() string {
	return "Users content"
}

func (this *Module) Module_users_sidebar() string {
	return "Users right sidebar"
}
