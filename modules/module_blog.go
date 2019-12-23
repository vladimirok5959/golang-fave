package modules

import (
	"html"
	"net/http"
	"strings"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/sqlw"
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
	}, func(wrap *wrapper.Wrapper) {
		if len(wrap.UrlArgs) == 3 && wrap.UrlArgs[0] == "blog" && wrap.UrlArgs[1] == "category" && wrap.UrlArgs[2] != "" {
			// Blog category
			row := &utils.MySql_blog_category{}
			rou := &utils.MySql_user{}
			err := wrap.DB.QueryRow(
				wrap.R.Context(),
				`SELECT
					main.id,
					main.user,
					main.name,
					main.alias,
					main.lft,
					main.rgt,
					main.depth,
					parent.id AS parent_id,
					users.id,
					users.first_name,
					users.last_name,
					users.email,
					users.admin,
					users.active
				FROM
					(
						SELECT
							node.id,
							node.user,
							node.name,
							node.alias,
							node.lft,
							node.rgt,
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
					) AS main
					LEFT JOIN (
						SELECT
							node.id,
							node.user,
							node.name,
							node.alias,
							node.lft,
							node.rgt,
							(COUNT(parent.id) - 0) AS depth
						FROM
							blog_cats AS node,
							blog_cats AS parent
						WHERE
							node.lft BETWEEN parent.lft AND parent.rgt
						GROUP BY
							node.id
						ORDER BY
							node.lft ASC
					) AS parent ON
					parent.depth = main.depth AND
					main.lft > parent.lft AND
					main.rgt < parent.rgt
					LEFT JOIN users ON users.id = main.user
				WHERE
					main.id > 1 AND
					main.alias = ?
				ORDER BY
					main.lft ASC
				;`,
				wrap.UrlArgs[2],
			).Scan(
				&row.A_id,
				&row.A_user,
				&row.A_name,
				&row.A_alias,
				&row.A_lft,
				&row.A_rgt,
				&row.A_depth,
				&row.A_parent,
				&rou.A_id,
				&rou.A_first_name,
				&rou.A_last_name,
				&rou.A_email,
				&rou.A_admin,
				&rou.A_active,
			)

			if err != nil && err != wrapper.ErrNoRows {
				// System error 500
				wrap.LogCpError(&err)
				utils.SystemErrorPageEngine(wrap.W, err)
				return
			} else if err == wrapper.ErrNoRows {
				// User error 404 page
				wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
				return
			}

			// Fix url
			if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
				http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
				return
			}

			// Render template
			wrap.RenderFrontEnd("blog-category", fetdata.New(wrap, false, row, rou), http.StatusOK)
			return
		} else if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "blog" && wrap.UrlArgs[1] != "" {
			// Blog post
			row := &utils.MySql_blog_post{}
			rou := &utils.MySql_user{}
			err := wrap.DB.QueryRow(
				wrap.R.Context(),
				`SELECT
					blog_posts.id,
					blog_posts.user,
					blog_posts.name,
					blog_posts.alias,
					blog_posts.category,
					blog_posts.briefly,
					blog_posts.content,
					UNIX_TIMESTAMP(blog_posts.datetime) as datetime,
					blog_posts.active,
					users.id,
					users.first_name,
					users.last_name,
					users.email,
					users.admin,
					users.active
				FROM
					blog_posts
					LEFT JOIN users ON users.id = blog_posts.user
				WHERE
					blog_posts.active = 1 and
					blog_posts.alias = ?
				LIMIT 1;`,
				wrap.UrlArgs[1],
			).Scan(
				&row.A_id,
				&row.A_user,
				&row.A_name,
				&row.A_alias,
				&row.A_category,
				&row.A_briefly,
				&row.A_content,
				&row.A_datetime,
				&row.A_active,
				&rou.A_id,
				&rou.A_first_name,
				&rou.A_last_name,
				&rou.A_email,
				&rou.A_admin,
				&rou.A_active,
			)

			if err != nil && err != wrapper.ErrNoRows {
				// System error 500
				wrap.LogCpError(&err)
				utils.SystemErrorPageEngine(wrap.W, err)
				return
			} else if err == wrapper.ErrNoRows {
				// User error 404 page
				wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
				return
			}

			// Fix url
			if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
				http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
				return
			}

			// Render template
			wrap.RenderFrontEnd("blog-post", fetdata.New(wrap, false, row, rou), http.StatusOK)
			return
		} else if len(wrap.UrlArgs) == 1 && wrap.UrlArgs[0] == "blog" {
			// Blog

			// Fix url
			if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
				http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
				return
			}

			// Render template
			wrap.RenderFrontEnd("blog", fetdata.New(wrap, false, nil, nil), http.StatusOK)
			return
		} else if (*wrap.Config).Engine.MainModule == 1 {
			// Render template
			wrap.RenderFrontEnd("blog", fetdata.New(wrap, false, nil, nil), http.StatusOK)
			return
		}

		// User error 404 page
		wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of posts"},
			})
			content += builder.DataTable(
				wrap,
				"blog_posts",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "name",
						NameInTable: "Post / URL",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							alias := html.EscapeString((*values)[2])
							return `<div>` + name + `</div><div><small>/blog/` + alias + `/</small></div>`
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
							Href:   `/blog/` + (*values)[2] + `/`,
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
							Href: "javascript:fave.ActionDataTableDelete(this,'blog-delete','" +
								(*values)[0] + "','Are you sure want to delete post?');",
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
					},
					{
						DBField: "user",
					},
					{
						DBField:     "name",
						NameInTable: "Category",
						CallBack: func(values *[]string) string {
							depth := utils.StrToInt((*values)[4]) - 1
							if depth < 0 {
								depth = 0
							}
							sub := strings.Repeat("&mdash; ", depth)
							name := `<a href="/cp/` + wrap.CurrModule + `/categories-modify/` + (*values)[0] + `/">` + sub + html.EscapeString((*values)[2]) + `</a>`
							return `<div>` + name + `</div>`
						},
					},
					{
						DBField: "alias",
					},
					{
						DBField: "depth",
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   `/blog/category/` + (*values)[3] + `/`,
							Hint:   "View",
							Target: "_blank",
						},
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
				nil,
				func(limit_offset int, pear_page int) (*sqlw.Rows, error) {
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
							node.lft BETWEEN parent.lft AND parent.rgt AND
							node.id > 1
						GROUP BY
							node.id
						ORDER BY
							node.lft ASC
						;`,
					)
				},
				false,
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

			data := utils.MySql_blog_post{
				A_id:       0,
				A_user:     0,
				A_name:     "",
				A_alias:    "",
				A_category: 0,
				A_content:  "",
				A_datetime: 0,
				A_active:   0,
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
						user,
						name,
						alias,
						category,
						briefly,
						content,
						active
					FROM
						blog_posts
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_name,
					&data.A_alias,
					&data.A_category,
					&data.A_briefly,
					&data.A_content,
					&data.A_active,
				)
				if *wrap.LogCpError(&err) != nil {
					return "", "", ""
				}
			}

			// All post current categories
			var selids []int
			if data.A_id > 0 {
				rows, err := wrap.DB.Query("SELECT category_id FROM blog_cat_post_rel WHERE post_id = ?;", data.A_id)
				if err == nil {
					defer rows.Close()
					values := make([]int, 1)
					scan := make([]interface{}, len(values))
					for i := range values {
						scan[i] = &values[i]
					}
					for rows.Next() {
						err = rows.Scan(scan...)
						if *wrap.LogCpError(&err) == nil {
							selids = append(selids, int(values[0]))
						}
					}
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
					Value: "blog-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Post name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Post alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: our-news",
					Max:     "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Category",
					Name:    "category",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n4">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_category">Category</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_category" name="category" data-live-search="true">` +
							`<option title="Nothing selected" value="0">&mdash;</option>` +
							this.blog_GetCategorySelectOptions(wrap, 0, data.A_category, []int{}) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKText,
					Caption: "Categories",
					Name:    "cats",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n5">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_parent">Categories</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_cats" name="cats[]" data-live-search="true" multiple>` +
							this.blog_GetCategorySelectOptions(wrap, 0, 0, selids) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Briefly",
					Name:    "briefly",
					Value:   data.A_briefly,
					Classes: "briefly wysiwyg",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Post content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "wysiwyg",
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
				err := wrap.DB.QueryRow(
					wrap.R.Context(),
					`SELECT
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
				if *wrap.LogCpError(&err) != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "categories-modify" {
				btn_caption = "Save"
			}

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
					Caption: "Parent",
					Name:    "parent",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n2">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_parent">Parent</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_parent" name="parent" data-live-search="true">` +
							`<option title="Nothing selected" value="0">&mdash;</option>` +
							this.blog_GetCategorySelectOptions(wrap, data.A_id, parentId, []int{}) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: popular-posts",
					Max:     "255",
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
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
