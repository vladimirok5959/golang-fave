package modules

import (
	utils "golang-fave/engine/wrapper/utils"
)

func (this *Module) Module_settings() {
	// Do something here...
}

func (this *Module) Module_settings_display() bool {
	return false
}

func (this *Module) Module_settings_alias() string {
	return "settings"
}

func (this *Module) Module_settings_name() string {
	return "Settings"
}

func (this *Module) Module_settings_submenu() []utils.ModuleSubMenu {
	result := make([]utils.ModuleSubMenu, 0)
	result = append(result, utils.ModuleSubMenu{Alias: "default", Name: "Settings"})
	return result
}

func (this *Module) Module_settings_content() string {
	return "Settings content"
}

func (this *Module) Module_settings_sidebar() string {
	return "Settings right sidebar"
}
