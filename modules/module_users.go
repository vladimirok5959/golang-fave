package modules

import (
	"html"

	"golang-fave/engine/assets"
	"golang-fave/engine/builder"
	"golang-fave/engine/consts"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_Users() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "users",
		Name:   "Users",
		Order:  800,
		System: true,
		Icon:   assets.SysSvgIconUser,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of users", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new user", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify user", Show: false},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of users"},
			})
			content += builder.DataTable(
				wrap,
				"fave_users",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "email",
						NameInTable: "Email / Name",
						CallBack: func(values *[]string) string {
							email := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							name := html.EscapeString((*values)[2])
							if name != "" && (*values)[3] != "" {
								name += ` ` + (*values)[3]
							}
							if name != "" {
								name = `<div><small>` + name + `</small></div>`
							}
							return `<div>` + email + `</div>` + name
						},
					},
					{
						DBField: "first_name",
					},
					{
						DBField: "last_name",
					},
					{
						DBField:     "active",
						NameInTable: "Active",
						Classes:     "d-none d-sm-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[4]))
						},
					},
					{
						DBField:     "admin",
						NameInTable: "Admin",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[5]))
						},
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'users-delete','" +
								(*values)[0] + "','Are you sure want to delete user?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/",
				nil,
				nil,
				true,
			)
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new user"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify user"},
				})
			}

			data := utils.MySql_user{
				A_id:         0,
				A_first_name: "",
				A_last_name:  "",
				A_email:      "",
				A_admin:      0,
				A_active:     0,
			}

			if wrap.CurrSubModule == "modify" {
				if len(wrap.UrlArgs) != 3 {
					return "", "", ""
				}
				if !utils.IsNumeric(wrap.UrlArgs[2]) {
					return "", "", ""
				}
				err := wrap.DB.QueryRow(
					wrap.R.Context(),
					`SELECT
						id,
						first_name,
						last_name,
						email,
						admin,
						active
					FROM
						fave_users
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_first_name,
					&data.A_last_name,
					&data.A_email,
					&data.A_admin,
					&data.A_active,
				)
				if *wrap.LogCpError(&err) != nil {
					return "", "", ""
				}
			}

			pass_req := true
			pass_hint := ""
			if wrap.CurrSubModule == "modify" {
				pass_req = false
				pass_hint = "Leave the field blank to not change the password"
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "users-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:    builder.DFKText,
					Caption: "First name",
					Name:    "first_name",
					Value:   data.A_first_name,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Last name",
					Name:    "last_name",
					Value:   data.A_last_name,
				},
				{
					Kind:     builder.DFKEmail,
					Caption:  "Email",
					Name:     "email",
					Value:    data.A_email,
					Required: true,
				},
				{
					Kind:     builder.DFKPassword,
					Caption:  "Password",
					Name:     "password",
					Required: pass_req,
					Hint:     pass_hint,
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Active",
					Name:    "active",
					Value:   utils.IntToStr(data.A_active),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Admin",
					Name:    "admin",
					Value:   utils.IntToStr(data.A_admin),
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  btn_caption,
					Target: "add-edit-button",
				},
			})

			if wrap.CurrSubModule == "add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
