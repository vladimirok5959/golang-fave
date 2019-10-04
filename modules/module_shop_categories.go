package modules

import (
	"html"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) shop_GetCategorySelectOptions(wrap *wrapper.Wrapper, id int, parentId int, selids []int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			node.id,
			node.user,
			node.name,
			node.alias,
			(COUNT(parent.id) - 1) AS depth
		FROM
			shop_cats AS node,
			shop_cats AS parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt AND
			node.id > 1
		GROUP BY
			node.id
		ORDER BY
			node.lft ASC
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 5)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		idStr := utils.IntToStr(id)
		parentIdStr := utils.IntToStr(parentId)
		for rows.Next() {
			err = rows.Scan(scan...)
			if *wrap.LogCpError(&err) == nil {
				disabled := ""
				if string(values[0]) == idStr {
					disabled = " disabled"
				}
				selected := ""
				if string(values[0]) == parentIdStr {
					selected = " selected"
				}
				if len(selids) > 0 && utils.InArrayInt(selids, utils.StrToInt(string(values[0]))) {
					selected = " selected"
				}
				depth := utils.StrToInt(string(values[4])) - 1
				if depth < 0 {
					depth = 0
				}
				sub := strings.Repeat("&mdash; ", depth)
				result += `<option title="` + html.EscapeString(string(values[2])) + `" value="` + html.EscapeString(string(values[0])) + `"` + disabled + selected + `>` + sub + html.EscapeString(string(values[2])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) shop_GetCategoryParentId(wrap *wrapper.Wrapper, id int) int {
	var parentId int
	err := wrap.DB.QueryRow(`
		SELECT
			parent.id
		FROM
			shop_cats AS node,
			shop_cats AS parent
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
	if *wrap.LogCpError(&err) != nil {
		return 0
	}
	return parentId
}
