package modules

import (
	"html"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) shop_GetProductsListForAttaching(wrap *wrapper.Wrapper, name string, id int) string {
	result := ``

	words := strings.Split(name, " ")
	filtered := []string{}
	for _, value := range words {
		word := utils.Trim(value)
		if word != "" {
			filtered = append(filtered, "%"+word+"%")
		}
	}

	search := ""
	params := make([]interface{}, len(filtered)+1)
	params[0] = id
	for i, value := range filtered {
		search += " AND name LIKE ?"
		params[i+1] = value
	}

	rows, err := wrap.DB.Query(
		wrap.R.Context(),
		`SELECT
			id,
			name
		FROM
			fave_shop_products
		WHERE
			id <> ? AND
			parent_id IS NULL
			`+search+`
		ORDER BY
			id DESC
		LIMIT 10;`,
		params...,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 2)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if *wrap.LogCpError(&err) == nil {
				result += `<div><a href="javascript:fave.ShopAttachProductTo(` + utils.IntToStr(id) + `,` + html.EscapeString(string(values[0])) + `);">` + html.EscapeString(string(values[1])) + ` ` + html.EscapeString(string(values[0])) + `</a></div>`
			}
		}
	}
	if result == "" {
		result = `<div><b>No any of products found</b></div>`
	}
	return result
}

func (this *Modules) RegisterAction_ShopAttachProductSearch() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-attach-product-search",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_words := utils.Trim(wrap.R.FormValue("words"))
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		if !utils.IsNumeric(pf_id) {
			wrap.Write(`$('#sys-modal-shop-product-attach .products-list').html('<b>Inner system error</b>');`)
			return
		}
		wrap.Write(`$('#sys-modal-shop-product-attach .products-list').html('` + this.shop_GetProductsListForAttaching(wrap, pf_words, utils.StrToInt(pf_id)) + `');`)
	})
}
