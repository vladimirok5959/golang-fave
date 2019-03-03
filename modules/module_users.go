package modules

import (
	"html"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
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
			content += builder.DataTable(wrap, "users", "id", "DESC", &[]builder.DataTableRow{
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
					CallBack: func(values *[]string) string {
						return builder.CheckBox(utils.StrToInt((*values)[4]))
					},
				},
				{
					DBField:     "admin",
					NameInTable: "Admin",
					CallBack: func(values *[]string) string {
						return builder.CheckBox(utils.StrToInt((*values)[5]))
					},
				},
			}, func(values *[]string) string {
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
						Hint: "Delete",
					},
				})
			}, "/cp/"+wrap.CurrModule+"/")
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
				err := wrap.DB.QueryRow(`
					SELECT
						id,
						first_name,
						last_name,
						email,
						admin,
						active
					FROM
						users
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
				if err != nil {
					return "", "", ""
				}
			}

			pass_req := true
			pass_hint := ""
			if wrap.CurrSubModule == "modify" {
				pass_req = false
				pass_hint = "Leave the field blank to not change the password"
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
					Kind: builder.DFKMessage,
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Add",
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

func (this *Modules) RegisterAction_UsersModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "users-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")
		pf_admin := wrap.R.FormValue("admin")
		pf_active := wrap.R.FormValue("active")

		if pf_admin == "" {
			pf_admin = "0"
		}

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		// First user always super admin
		// Rewrite active and admin status
		if pf_id == "1" {
			pf_admin = "1"
			pf_active = "1"
		}

		if pf_id == "0" {
			// Add new user
			if pf_password == "" {
				wrap.MsgError(`Please specify user password`)
				return
			}
			_, err := wrap.DB.Query(
				`INSERT INTO users SET
					first_name = ?,
					last_name = ?,
					email = ?,
					password = MD5(?),
					admin = ?,
					active = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				pf_password,
				pf_admin,
				pf_active,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/users/';`)
		} else {
			// Update user
			if pf_password == "" {
				_, err := wrap.DB.Query(
					`UPDATE users SET
						first_name = ?,
						last_name = ?,
						email = ?,
						admin = ?,
						active = ?
					WHERE
						id = ?
					;`,
					pf_first_name,
					pf_last_name,
					pf_email,
					pf_admin,
					pf_active,
					utils.StrToInt(pf_id),
				)
				if err != nil {
					wrap.MsgError(err.Error())
					return
				}
			} else {
				_, err := wrap.DB.Query(
					`UPDATE users SET
						first_name = ?,
						last_name = ?,
						email = ?,
						password = MD5(?)
					WHERE
						id = ?
					;`,
					pf_first_name,
					pf_last_name,
					pf_email,
					pf_password,
					utils.StrToInt(pf_id),
				)
				if err != nil {
					wrap.MsgError(err.Error())
					return
				}
			}
			wrap.Write(`window.location='/cp/users/modify/` + pf_id + `/';`)
		}
	})
}

func (this *Modules) RegisterAction_UsersDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "users-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Delete user
		_, err := wrap.DB.Query(
			`DELETE FROM users WHERE id = ? and id <> 1;`,
			utils.StrToInt(pf_id),
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
