package modules

import (
	"html"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopGetAttributeValues() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-get-attribute-values",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		options := ``
		rows, err := wrap.DB.Query(
			`SELECT
				id,
				name
			FROM
				shop_filters_values
			WHERE
				filter_id = ?
			ORDER BY
				name ASC
			;`,
			utils.StrToInt(pf_id),
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
				if err == nil {
					options += `<option value="` + html.EscapeString(string(values[0])) + `">` + html.EscapeString(utils.JavaScriptVarValue(string(values[1]))) + `</option>`
				}
			}
		}

		wrap.Write(`if($('#prod_attr_` + pf_id + `').length > 0) {`)
		wrap.Write(`$('#prod_attr_` + pf_id + ` select').prop('disabled', false).prop('multiple', true);`)
		wrap.Write(`$('#prod_attr_` + pf_id + ` select').html('` + options + `');`)
		wrap.Write(`$('#prod_attr_` + pf_id + ` select').selectpicker({});`)
		wrap.Write(`$('#prod_attr_` + pf_id + ` button').prop('disabled', false);`)
		wrap.Write(`}`)
	})
}
