package modules

import (
	"html"

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
	if this.smod == "default" {
		breadcrumb := this.breadcrumb([]dataBreadcrumb{
			{"List of users", ""},
		})
		data_table := this.data_table("users", "email", "ASC", []dataTableRow{
			{"id", "", nil},
			{"email", "Email", func(values *[]string) string {
				email := `<a href="/cp/users/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
				name := html.EscapeString((*values)[2])
				if name != "" && (*values)[3] != "" {
					name += ` ` + (*values)[3]
				}
				if name != "" {
					name = `<div><small>` + name + `</small></div>`
				}
				return `<div>` + email + `</div>` + name
			}},
			{"first_name", "", nil},
			{"last_name", "", nil},
		}, func(values *[]string) string {
			return `<a href="/cp/users/modify/` + (*values)[0] + `/">` +
				others.File_assets_sys_svg_edit + `</a>` +
				`<a href="#">` + others.File_assets_sys_svg_remove + `</a>`
		}, "/cp/users/default/")
		return breadcrumb + data_table
	} else if this.smod == "modify" {
		// Add/Edit
		return "Edit!"
	}
	return ""
}

func (this *Module) Module_users_sidebar() string {
	return ""
}
