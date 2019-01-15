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
			return `<a class="ico" href="/cp/users/modify/` + (*values)[0] + `/">` +
				others.File_assets_sys_svg_edit + `</a>` +
				`<a class="ico" href="#">` + others.File_assets_sys_svg_remove + `</a>`
		}, "/cp/users/default/")
		return breadcrumb + data_table
	} else if this.smod == "modify" && this.imod == 0 {
		breadcrumb := this.breadcrumb([]dataBreadcrumb{
			{"Add new user", ""},
		})
		data_form := this.data_form([]dataFormField{
			{kind: dfkHidden, name: "action", value: "users_add"},
			{kind: dfkHidden, name: "id", value: "0"},
			{kind: dfkText, caption: "User first name", name: "first_name"},
			{kind: dfkText, caption: "User last name", name: "last_name"},
			{kind: dfkEmail, caption: "User email", name: "email", required: true},
			{kind: dfkPassword, caption: "User password", name: "password", hint: "Please specify new user password", required: true},
			{kind: dfkSubmit, value: "Add", target: "add-edit-button"},
		})
		return breadcrumb + data_form
	} else if this.smod == "modify" && this.imod != 0 {
		breadcrumb := this.breadcrumb([]dataBreadcrumb{
			{"Edit user", ""},
		})
		// Load user data
		data_form := this.data_form([]dataFormField{
			{kind: dfkHidden, name: "action", value: "users_edit"},
			{kind: dfkHidden, name: "id", value: "0"},
			{kind: dfkText, caption: "User first name", name: "first_name", value: "1"},
			{kind: dfkText, caption: "User last name", name: "last_name", value: "2"},
			{kind: dfkEmail, caption: "User email", name: "email", value: "3", required: true},
			{kind: dfkPassword, caption: "User password", name: "password", hint: "Leave the field blank to not change the password"},
			{kind: dfkSubmit, value: "Add", target: "add-edit-button"},
		})
		return breadcrumb + data_form
	}
	return ""
}

func (this *Module) Module_users_sidebar() string {
	if this.smod == "modify" && this.imod == 0 {
		return `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
	} else if this.smod == "modify" && this.imod != 0 {
		return `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
	}
	return ""
}
