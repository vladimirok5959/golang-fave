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
	if this.smod == "default" {
		// List
		result := `<table class="table table-striped table-bordered">
			<thead>
				<tr>
					<th scope="col">Email</th>
					<th scope="col">First name</th>
					<th scope="col">Last name</th>
					<th scope="col">Action</th>
				</tr>
			</thead>
		<tbody>`
		rows, err := this.db.Query("SELECT `id`, `first_name`, `last_name`, `email` FROM `users`;")
		if err == nil {
			var id int
			var first_name string
			var last_name string
			var email string
			for rows.Next() {
				err = rows.Scan(&id, &first_name, &last_name, &email)
				if err == nil {
					result += `<tr>
						<td>` + email + `</td>
						<td>` + first_name + `</td>
						<td>` + last_name + `</td>
						<td><a href="#">` + others.File_assets_sys_svg_edit + `</a> <a href="#">` + others.File_assets_sys_svg_remove + `</a></td>
					</tr>`
				}
			}
		}
		result += `</tbody></table>`
		return result
	} else if this.smod == "modify" {
		// Add/Edit
	}
	return ""
}

func (this *Module) Module_users_sidebar() string {
	return ""
}
