package modules

import (
	"html"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) blog_GetCategorySelectOptions(wrap *wrapper.Wrapper, id int, parentId int) string {
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
		parentIdStr := utils.IntToStr(parentId)
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				disabled := ""
				if string(values[0]) == idStr {
					disabled = " disabled"
				}
				selected := ""
				if string(values[0]) == parentIdStr {
					selected = " selected"
				}
				sub := strings.Repeat("&mdash; ", utils.StrToInt(string(values[4])))
				result += `<option value="` + html.EscapeString(string(values[0])) + `"` + disabled + selected + `>` + sub + html.EscapeString(string(values[2])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) blog_GetCategoryParentId(wrap *wrapper.Wrapper, id int) int {
	var parentId int
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
		&parentId,
	)
	return parentId
}

func (this *Modules) blog_ActionCategoryAdd(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	//
	return nil
}

func (this *Modules) blog_ActionCategoryUpdate(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	//
	return nil
}

func (this *Modules) RegisterAction_BlogCategoriesModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-categories-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_parent := wrap.R.FormValue("parent")

		if !utils.IsNumeric(pf_id) || !utils.IsNumeric(pf_parent) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify category name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateSingleAlias(pf_name)
		}

		if !utils.IsValidSingleAlias(pf_alias) {
			wrap.MsgError(`Please specify correct category alias`)
			return
		}

		if pf_id == "0" {
			if err := this.blog_ActionCategoryAdd(wrap, pf_id, pf_name, pf_alias, pf_parent); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/blog/categories/';`)
		} else {
			if err := this.blog_ActionCategoryUpdate(wrap, pf_id, pf_name, pf_alias, pf_parent); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/blog/categories-modify/` + pf_id + `/';`)
		}
	})
}
