package modules

import (
	"database/sql"
	"html"
	"strings"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_Blog() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "blog",
		Name:   "Blog",
		Order:  1,
		System: false,
		Icon:   assets.SysSvgIconList,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of posts", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new post", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify post", Show: false},
			{Sep: true, Show: true},
			{Mount: "categories", Name: "List of categories", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "categories-add", Name: "Add new category", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "categories-modify", Name: "Modify category", Show: false},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of posts"},
			})
			//
		} else if wrap.CurrSubModule == "categories" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of categories"},
			})
			content += builder.DataTable(
				wrap,
				"blog_cats",
				"id",
				"ASC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
						// NameInTable: "id",
					},
					{
						DBField: "user",
						// NameInTable: "user",
					},
					{
						DBField:     "name",
						NameInTable: "Category",
						CallBack: func(values *[]string) string {
							sub := strings.Repeat("&mdash; ", utils.StrToInt((*values)[4]))
							name := `<a href="/cp/` + wrap.CurrModule + `/categories-modify/` + (*values)[0] + `/">` + sub + html.EscapeString((*values)[2]) + `</a>`
							// alias := html.EscapeString((*values)[3])
							// return `<div>` + name + `</div><div><small>` + alias + `</small></div>`
							return `<div>` + name + `</div>`
						},
					},
					{
						DBField: "alias",
						// NameInTable: "Alias",
					},
					{
						DBField: "depth",
						// NameInTable: "depth",
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/categories-modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'blog-categories-delete','" +
								(*values)[0] + "','Are you sure want to delete category?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/"+wrap.CurrSubModule+"/",
				func() (int, error) {
					var num int
					var err error
					err = wrap.DB.QueryRow("SELECT COUNT(*) FROM `blog_cats`;").Scan(&num)
					return num, err
				},
				func(limit_offset int, pear_page int) (*sql.Rows, error) {
					return wrap.DB.Query(
						`SELECT
							node.id,
							node.user,
							node.name,
							node.alias,
							(COUNT(parent.id) - 1) AS depth
						FROM
							blog_cats AS node,
							blog_cats AS parent
						WHERE
							node.lft BETWEEN parent.lft AND parent.rgt
						GROUP BY
							node.id
						ORDER BY
							node.lft ASC
						LIMIT ?, ?;`,
						limit_offset,
						pear_page,
					)
				},
			)
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new post"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify post"},
				})
			}
			//
		} else if wrap.CurrSubModule == "categories-add" || wrap.CurrSubModule == "categories-modify" {
			if wrap.CurrSubModule == "categories-add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new category"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify category"},
				})
			}
			// ---

			btn_caption := "Add"
			if wrap.CurrSubModule == "categories-modify" {
				btn_caption = "Save"
			}

			// ---
			select_parent_options := ""
			rows, err := wrap.DB.Query(
				`SELECT
					node.id,
					node.user,
					node.name,
					node.alias,
					(COUNT(parent.id) - 1) AS depth
				FROM
					blog_cats AS node,
					blog_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt
				GROUP BY
					node.id
				ORDER BY
					node.lft ASC
				;`,
			)
			if err == nil {
				values := make([]string, 5)
				scan := make([]interface{}, len(values))
				for i := range values {
					scan[i] = &values[i]
				}
				for rows.Next() {
					err = rows.Scan(scan...)
					if err == nil {
						sub := strings.Repeat("&mdash; ", utils.StrToInt(string(values[4])))
						select_parent_options += `<option value="` + html.EscapeString(string(values[0])) + `">` + sub + html.EscapeString(string(values[2])) + `</option>`
					}
				}
			}
			// ---

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "blog-categories-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: "0",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Name",
					Name:    "name",
					Value:   "",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Alias",
					Name:    "alias",
					Value:   "",
					Hint:    "Example: popular-posts",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Parent",
					Name:    "parent",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n3">
							<div class="row">
								<div class="col-md-3">
									<label for="lbl_parent">Parent</label>
								</div>
								<div class="col-md-9">
									<div>
										<select class="form-control" id="lbl_parent" name="parent">
											<option value="0">&mdash;</option>
											` + select_parent_options + `
										</select>
									</div>
								</div>
							</div>
						</div>`
					},
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

			if wrap.CurrSubModule == "categories-add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
			// ---
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
