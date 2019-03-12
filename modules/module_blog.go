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

func (this *Modules) blog_GetCategorySelectOptions(wrap *wrapper.Wrapper, id int) string {
	result := ``
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
		idStr := utils.IntToStr(id)
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				selected := ""
				if string(values[0]) == idStr {
					selected = " selected"
				}
				sub := strings.Repeat("&mdash; ", utils.StrToInt(string(values[4])))
				result += `<option value="` + html.EscapeString(string(values[0])) + `"` + selected + `>` + sub + html.EscapeString(string(values[2])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) blog_GetCategoryParentId(wrap *wrapper.Wrapper, id int) int {
	var result int
	_ = wrap.DB.QueryRow(`
		SELECT
			parent.id
		FROM
			blog_cats AS node,
			blog_cats AS parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt AND
			node.id = ? AND
			parent.id <> ?
		ORDER BY
			parent.lft DESC
		LIMIT 1;`,
		id,
		id,
	).Scan(
		&result,
	)
	return result
}

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
				{Name: "Categories", Link: "/cp/" + wrap.CurrModule + "/" + wrap.CurrSubModule + "/"},
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
					{Name: "Categories", Link: "/cp/" + wrap.CurrModule + "/categories/"},
					{Name: "Add new category"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Categories", Link: "/cp/" + wrap.CurrModule + "/categories/"},
					{Name: "Modify category"},
				})
			}

			data := utils.MySql_blog_category{
				A_id:    0,
				A_user:  0,
				A_name:  "",
				A_alias: "",
				A_lft:   0,
				A_rgt:   0,
			}

			if wrap.CurrSubModule == "categories-modify" {
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
						lft,
						rgt
					FROM
						blog_cats
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_name,
					&data.A_alias,
					&data.A_lft,
					&data.A_rgt,
				)
				if err != nil {
					return "", "", ""
				}
			}

			// ---

			btn_caption := "Add"
			if wrap.CurrSubModule == "categories-modify" {
				btn_caption = "Save"
			}

			// ---
			// select_parent_options := this.blog_GetCategorySelectOptions(wrap, "")
			// ---
			parentId := 0
			if wrap.CurrSubModule == "categories-modify" {
				parentId = this.blog_GetCategoryParentId(wrap, data.A_id)
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "blog-categories-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:    builder.DFKText,
					Caption: "Name",
					Name:    "name",
					Value:   data.A_name,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: popular-posts",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Parent",
					Name:    "parent",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group last">
							<div class="row">
								<div class="col-md-3">
									<label for="lbl_parent">Parent</label>
								</div>
								<div class="col-md-9">
									<div>
										<select class="form-control" id="lbl_parent" name="parent">
											<option value="0">&mdash;</option>
											` + this.blog_GetCategorySelectOptions(wrap, parentId) + `
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
