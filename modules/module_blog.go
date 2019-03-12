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
						NameInTable: "Name",
						CallBack: func(values *[]string) string {
							sub := strings.Repeat("— ", utils.StrToInt((*values)[4]))
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
			//
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
