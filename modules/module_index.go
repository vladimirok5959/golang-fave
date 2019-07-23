package modules

import (
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
		if err != nil && err != wrapper.ErrNoRows {
			// System error 500
			utils.SystemErrorPageEngine(wrap.W, err)
			return
		} else if err == wrapper.ErrNoRows {
			// User error 404 page
			wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true), http.StatusNotFound)
			return
		}

		// Which template
		tname := "index"
		if wrap.R.URL.Path != "/" {
			tname = "page"
		}

		// Render template
		wrap.RenderFrontEnd(tname, fetdata.New(wrap, row, false), http.StatusOK)
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
						NameInTable: "Page / URL",
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
				true,
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
					Kind:     builder.DFKText,
					Caption:  "Page name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: /about-us/ or /about-us.html",
					Max:     "255",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Page content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "wysiwyg",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta title",
					Name:    "meta_title",
					Value:   data.A_meta_title,
					Max:     "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta keywords",
					Name:    "meta_keywords",
					Value:   data.A_meta_keywords,
					Max:     "255",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Meta description",
					Name:    "meta_description",
					Value:   data.A_meta_description,
					Max:     "510",
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Active",
					Name:    "active",
					Value:   utils.IntToStr(data.A_active),
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
