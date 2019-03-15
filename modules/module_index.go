package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"html"
	"net/http"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_Index() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "index",
		Name:   "Pages",
		Order:  0,
		Icon:   assets.SysSvgIconPage,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of pages", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new page", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify page", Show: false},
		},
	}, func(wrap *wrapper.Wrapper) {
		// Front-end
		row := &utils.MySql_page{}
		err := wrap.DB.QueryRow(`
			SELECT
				id,
				user,
				name,
				alias,
				content,
				meta_title,
				meta_keywords,
				meta_description,
				UNIX_TIMESTAMP(datetime) as datetime,
				active
			FROM
				pages
			WHERE
				active = 1 and
				alias = ?
			LIMIT 1;`,
			wrap.R.URL.Path,
		).Scan(
			&row.A_id,
			&row.A_user,
			&row.A_name,
			&row.A_alias,
			&row.A_content,
			&row.A_meta_title,
			&row.A_meta_keywords,
			&row.A_meta_description,
			&row.A_datetime,
			&row.A_active,
		)
		if err != nil && err != sql.ErrNoRows {
			// System error 500
			utils.SystemErrorPageEngine(wrap.W, err)
			return
		} else if err == sql.ErrNoRows {
			// User error 404 page
			wrap.W.WriteHeader(http.StatusNotFound)
			wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true))
			return
		}

		// Replace title with page name
		if row.A_meta_title == "" {
			row.A_meta_title = row.A_name
		}

		// Which template
		tmpl_name := "index"
		if wrap.R.URL.Path != "/" {
			tmpl_name = "page"
		}

		// Render template
		wrap.RenderFrontEnd(tmpl_name, fetdata.New(wrap, row, false))
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of pages"},
			})
			content += builder.DataTable(
				wrap,
				"pages",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "name",
						NameInTable: "Page / Alias",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							alias := html.EscapeString((*values)[2])
							return `<div>` + name + `</div><div><small>` + alias + `</small></div>`
						},
					},
					{
						DBField: "alias",
					},
					{
						DBField:     "datetime",
						DBExp:       "UNIX_TIMESTAMP(`datetime`)",
						NameInTable: "Date / Time",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							t := int64(utils.StrToInt((*values)[3]))
							return `<div>` + utils.UnixTimestampToFormat(t, "02.01.2006") + `</div>` +
								`<div><small>` + utils.UnixTimestampToFormat(t, "15:04:05") + `</small></div>`
						},
					},
					{
						DBField:     "active",
						NameInTable: "Active",
						Classes:     "d-none d-sm-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[4]))
						},
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   (*values)[2],
							Hint:   "View",
							Target: "_blank",
						},
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'index-delete','" +
								(*values)[0] + "','Are you sure want to delete page?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/",
				nil,
				nil,
			)
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new page"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify page"},
				})
			}

			data := utils.MySql_page{
				A_id:               0,
				A_user:             0,
				A_name:             "",
				A_alias:            "",
				A_content:          "",
				A_meta_title:       "",
				A_meta_keywords:    "",
				A_meta_description: "",
				A_datetime:         0,
				A_active:           0,
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
						user,
						name,
						alias,
						content,
						meta_title,
						meta_keywords,
						meta_description,
						active
					FROM
						pages
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_name,
					&data.A_alias,
					&data.A_content,
					&data.A_meta_title,
					&data.A_meta_keywords,
					&data.A_meta_description,
					&data.A_active,
				)
				if err != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "index-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page name",
					Name:    "name",
					Value:   data.A_name,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: /about-us/ or /about-us.html",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Page content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "autosize",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta title",
					Name:    "meta_title",
					Value:   data.A_meta_title,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta keywords",
					Name:    "meta_keywords",
					Value:   data.A_meta_keywords,
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Meta description",
					Name:    "meta_description",
					Value:   data.A_meta_description,
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Active",
					Name:    "active",
					Value:   utils.IntToStr(data.A_active),
				},
				{
					Kind: builder.DFKMessage,
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

func (this *Modules) RegisterAction_IndexModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "index-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_content := wrap.R.FormValue("content")
		pf_meta_title := wrap.R.FormValue("meta_title")
		pf_meta_keywords := wrap.R.FormValue("meta_keywords")
		pf_meta_description := wrap.R.FormValue("meta_description")
		pf_active := wrap.R.FormValue("active")

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify page name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateAlias(pf_name)
		}

		if !utils.IsValidAlias(pf_alias) {
			wrap.MsgError(`Please specify correct page alias`)
			return
		}

		if pf_id == "0" {
			// Add new page
			_, err := wrap.DB.Query(
				`INSERT INTO pages SET
					user = ?,
					name = ?,
					alias = ?,
					content = ?,
					meta_title = ?,
					meta_keywords = ?,
					meta_description = ?,
					datetime = ?,
					active = ?
				;`,
				wrap.User.A_id,
				pf_name,
				pf_alias,
				pf_content,
				pf_meta_title,
				pf_meta_keywords,
				pf_meta_description,
				utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
				pf_active,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/';`)
		} else {
			// Update page
			_, err := wrap.DB.Query(
				`UPDATE pages SET
					name = ?,
					alias = ?,
					content = ?,
					meta_title = ?,
					meta_keywords = ?,
					meta_description = ?,
					active = ?
				WHERE
					id = ?
				;`,
				pf_name,
				pf_alias,
				pf_content,
				pf_meta_title,
				pf_meta_keywords,
				pf_meta_description,
				pf_active,
				utils.StrToInt(pf_id),
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/index/modify/` + pf_id + `/';`)
		}
	})
}

func (this *Modules) RegisterAction_IndexDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "index-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Delete page
		_, err := wrap.DB.Query(
			`DELETE FROM pages WHERE id = ?;`,
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

func (this *Modules) RegisterAction_IndexUserUpdateProfile() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "index-user-update-profile",
		WantUser: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		if pf_password != "" {
			// Update with password if set
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
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		} else {
			// Update without password if not set
			_, err := wrap.DB.Query(
				`UPDATE users SET
					first_name = ?,
					last_name = ?,
					email = ?
				WHERE
					id = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
