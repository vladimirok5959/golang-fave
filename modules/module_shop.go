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

func (this *Modules) shop_GetCurrencySelectOptions(wrap *wrapper.Wrapper, id int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			id,
			code
		FROM
			shop_currencies
		ORDER BY
			id ASC
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 2)
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
				result += `<option title="` + html.EscapeString(string(values[1])) + `" value="` + html.EscapeString(string(values[0])) + `"` + selected + `>` + html.EscapeString(string(values[1])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) shop_GetProductValuesInputs(wrap *wrapper.Wrapper, product_id int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			shop_filters.id,
			shop_filters.name,
			shop_filters_values.id,
			shop_filters_values.name,
			IF(shop_filter_product_values.filter_value_id > 0, 1, 0) as selected
		FROM
			shop_filters_values
			LEFT JOIN shop_filters ON shop_filters.id = shop_filters_values.filter_id
			LEFT JOIN shop_filter_product_values ON
				shop_filter_product_values.filter_value_id = shop_filters_values.id AND
				shop_filter_product_values.product_id = ` + utils.IntToStr(product_id) + `
			LEFT JOIN (
				SELECT
					shop_filters_values.filter_id,
					shop_filter_product_values.product_id
				FROM
					shop_filter_product_values
					LEFT JOIN shop_filters_values ON shop_filters_values.id = shop_filter_product_values.filter_value_id 
				WHERE
					shop_filter_product_values.product_id = ` + utils.IntToStr(product_id) + `
				GROUP BY
					shop_filters_values.filter_id
			) as filter_used ON filter_used.filter_id = shop_filters.id
		WHERE
			filter_used.filter_id IS NOT NULL
		ORDER BY
			shop_filters.name ASC,
			shop_filters_values.name ASC
		;`,
	)

	filter_ids := []int{}
	filter_names := map[int]string{}
	filter_values := map[int][]string{}

	if err == nil {
		defer rows.Close()
		values := make([]string, 5)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				filter_id := utils.StrToInt(string(values[0]))
				if !utils.InArrayInt(filter_ids, filter_id) {
					filter_ids = append(filter_ids, filter_id)
				}
				filter_names[filter_id] = html.EscapeString(string(values[1]))
				selected := ``
				if utils.StrToInt(string(values[4])) == 1 {
					selected = ` selected`
				}
				filter_values[filter_id] = append(filter_values[filter_id], `<option value="`+html.EscapeString(string(values[2]))+`"`+selected+`>`+html.EscapeString(string(values[3]))+`</option>`)
			}
		}
	}
	for _, filter_id := range filter_ids {
		result += `<div class="form-group" id="prod_attr_` + utils.IntToStr(filter_id) + `">` +
			`<div><b>` + filter_names[filter_id] + `</b></div>` +
			`<div class="position-relative">` +
			`<select class="selectpicker form-control" name="value.` + utils.IntToStr(filter_id) + `" autocomplete="off" required multiple>` +
			strings.Join(filter_values[filter_id], "") +
			`</select>` +
			`<button type="button" class="btn btn-danger btn-dynamic-remove" onclick="fave.ShopProductsRemove(this);">&times;</button>` +
			`</div>` +
			`</div>`
	}
	return result
}

func (this *Modules) shop_GetFilterValuesInputs(wrap *wrapper.Wrapper, filter_id int) string {
	result := ``
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
		filter_id,
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
				result += `<div class="form-group position-relative"><input class="form-control" type="text" name="value.` + html.EscapeString(string(values[0])) + `" value="` + html.EscapeString(string(values[1])) + `" placeholder="" autocomplete="off" required><button type="button" class="btn btn-danger btn-dynamic-remove" onclick="fave.ShopAttributesRemove(this);">&times;</button></div>`
			}
		}
	}
	return result
}

func (this *Modules) shop_GetAllAttributesSelectOptions(wrap *wrapper.Wrapper) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			id,
			name,
			filter
		FROM
			shop_filters
		ORDER BY
			name ASC
		;`,
	)
	result += `<option title="&mdash;" value="0">&mdash;</option>`
	if err == nil {
		defer rows.Close()
		values := make([]string, 3)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<option title="` + html.EscapeString(string(values[1])) + `" value="` + html.EscapeString(string(values[0])) + `">` + html.EscapeString(string(values[1])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) shop_GetAllCurrencies(wrap *wrapper.Wrapper) map[int]string {
	result := map[int]string{}
	rows, err := wrap.DB.Query(
		`SELECT
			id,
			code
		FROM
			shop_currencies
		ORDER BY
			id ASC
		;`,
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
				result[utils.StrToInt(string(values[0]))] = html.EscapeString(string(values[1]))
			}
		}
	}
	return result
}

func (this *Modules) shop_GetAllProductImages(wrap *wrapper.Wrapper, product_id int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			id,
			product_id,
			filename
		FROM
			shop_product_images
		WHERE
			product_id = ?
		ORDER BY
			ord ASC
		;`,
		product_id,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 3)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<div class="attached-img" data-id="` + html.EscapeString(string(values[0])) + `"><a href="/products/images/` + html.EscapeString(string(values[1])) + `/` + html.EscapeString(string(values[2])) + `" title="` + html.EscapeString(string(values[2])) + `" target="_blank"><img id="pimg_` + string(values[1]) + `_` + strings.Replace(string(values[2]), ".", "_", -1) + `" src="/products/images/` + string(values[1]) + `/thumb-0-` + string(values[2]) + `" onerror="WaitForFave(function(){fave.ShopProductsRetryImage(this, 'pimg_` + string(values[1]) + `_` + strings.Replace(string(values[2]), ".", "_", -1) + `');});" /></a><a class="remove" onclick="fave.ShopProductsDeleteImage(this, ` + html.EscapeString(string(values[1])) + `, '` + html.EscapeString(string(values[2])) + `');"><svg viewBox="1 1 11 14" width="10" height="12" class="sicon" version="1.1"><path fill-rule="evenodd" d="M11 2H9c0-.55-.45-1-1-1H5c-.55 0-1 .45-1 1H2c-.55 0-1 .45-1 1v1c0 .55.45 1 1 1v9c0 .55.45 1 1 1h7c.55 0 1-.45 1-1V5c.55 0 1-.45 1-1V3c0-.55-.45-1-1-1zm-1 12H3V5h1v8h1V5h1v8h1V5h1v8h1V5h1v9zm1-10H2V3h9v1z"></path></svg></a></div>`
			}
		}
	}
	return result
}

func (this *Modules) RegisterModule_Shop() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "shop",
		Name:   "Shop",
		Order:  2,
		System: false,
		Icon:   assets.SysSvgIconShop,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of products", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new product", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify product", Show: false},
			{Sep: true, Show: true},
			{Mount: "categories", Name: "List of categories", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "categories-add", Name: "Add new category", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "categories-modify", Name: "Modify category", Show: false},
			{Sep: true, Show: true},
			{Mount: "attributes", Name: "List of attributes", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "attributes-add", Name: "Add new attribute", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "attributes-modify", Name: "Modify attribute", Show: false},
			{Sep: true, Show: true},
			{Mount: "currencies", Name: "List of currencies", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "currencies-add", Name: "Add new currency", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "currencies-modify", Name: "Modify currency", Show: false},
		},
	}, func(wrap *wrapper.Wrapper) {
		if len(wrap.UrlArgs) == 3 && wrap.UrlArgs[0] == "shop" && wrap.UrlArgs[1] == "category" && wrap.UrlArgs[2] != "" {
			// Shop category
			row := &utils.MySql_shop_category{}
			rou := &utils.MySql_user{}
			err := wrap.DB.QueryRow(`
				SELECT
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
							shop_cats AS node,
							shop_cats AS parent
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
							shop_cats AS node,
							shop_cats AS parent
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
			wrap.RenderFrontEnd("shop-category", fetdata.New(wrap, false, row, rou), http.StatusOK)
			return
		} else if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "shop" && wrap.UrlArgs[1] != "" {
			// Shop product
			row := &utils.MySql_shop_product{}
			rou := &utils.MySql_user{}
			err := wrap.DB.QueryRow(`
				SELECT
					shop_products.id,
					shop_products.user,
					shop_products.currency,
					shop_products.price,
					shop_products.name,
					shop_products.alias,
					shop_products.vendor,
					shop_products.quantity,
					shop_products.category,
					shop_products.briefly,
					shop_products.content,
					UNIX_TIMESTAMP(shop_products.datetime) as datetime,
					shop_products.active,
					users.id,
					users.first_name,
					users.last_name,
					users.email,
					users.admin,
					users.active
				FROM
					shop_products
					LEFT JOIN users ON users.id = shop_products.user
				WHERE
					shop_products.active = 1 and
					shop_products.alias = ?
				LIMIT 1;`,
				wrap.UrlArgs[1],
			).Scan(
				&row.A_id,
				&row.A_user,
				&row.A_currency,
				&row.A_price,
				&row.A_name,
				&row.A_alias,
				&row.A_vendor,
				&row.A_quantity,
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
			wrap.RenderFrontEnd("shop-product", fetdata.New(wrap, false, row, rou), http.StatusOK)
			return
		} else if len(wrap.UrlArgs) == 1 && wrap.UrlArgs[0] == "shop" {
			// Shop

			// Fix url
			if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
				http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
				return
			}

			// Render template
			wrap.RenderFrontEnd("shop", fetdata.New(wrap, false, nil, nil), http.StatusOK)
			return
		} else if (*wrap.Config).Engine.MainModule == 2 {
			// Render template
			wrap.RenderFrontEnd("shop", fetdata.New(wrap, false, nil, nil), http.StatusOK)
			return
		}

		// User error 404 page
		wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of products"},
			})

			// Load currencies
			currencies := this.shop_GetAllCurrencies(wrap)

			content += builder.DataTable(
				wrap,
				"shop_products",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "name",
						NameInTable: "Product / URL",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							alias := html.EscapeString((*values)[2])
							parent := ``
							if (*values)[7] != "" {
								parent = `<div>&uarr;<small><a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[7] + `/">` + html.EscapeString((*values)[8]) + `</a></small></div>`
							}
							return `<div>` + name + `</div><div><small>/shop/` + alias + `/</small></div>` + parent
						},
					},
					{
						DBField: "alias",
					},
					{
						DBField: "currency",
					},
					{
						DBField:     "price",
						NameInTable: "Price",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							return `<div>` + utils.Float64ToStr(utils.StrToFloat64((*values)[4])) + `</div>` +
								`<div><small>` + currencies[utils.StrToInt((*values)[3])] + `</small></div>`
						},
					},
					{
						DBField:     "datetime",
						DBExp:       "UNIX_TIMESTAMP(`datetime`)",
						NameInTable: "Date / Time",
						Classes:     "d-none d-lg-table-cell",
						CallBack: func(values *[]string) string {
							t := int64(utils.StrToInt((*values)[5]))
							return `<div>` + utils.UnixTimestampToFormat(t, "02.01.2006") + `</div>` +
								`<div><small>` + utils.UnixTimestampToFormat(t, "15:04:05") + `</small></div>`
						},
					},
					{
						DBField:     "active",
						NameInTable: "Active",
						Classes:     "d-none d-sm-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[6]))
						},
					},
					{
						DBField: "parent_id",
					},
					{
						DBField: "pname",
						DBExp:   "spp.name",
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   `/shop/` + (*values)[2] + `/`,
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
							Href: "javascript:fave.ActionDataTableDelete(this,'shop-delete','" +
								(*values)[0] + "','Are you sure want to delete product?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/",
				func() (int, error) {
					var count int
					return count, wrap.DB.QueryRow(
						"SELECT COUNT(*) FROM `shop_products`;",
					).Scan(&count)
				},
				func(limit_offset int, pear_page int) (*sqlw.Rows, error) {
					return wrap.DB.Query(
						`SELECT
							shop_products.id,
							shop_products.name,
							shop_products.alias,
							shop_products.currency,
							shop_products.price,
							UNIX_TIMESTAMP(`+"`shop_products`.`datetime`"+`) AS datetime,
							shop_products.active,
							shop_products.parent_id,
							spp.name AS pname
						FROM
							shop_products
							LEFT JOIN shop_products AS spp ON spp.id = shop_products.parent_id
						ORDER BY
							shop_products.id DESC
						LIMIT ?, ?;`,
						limit_offset,
						pear_page,
					)
				},
				true,
			)
		} else if wrap.CurrSubModule == "categories" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Categories", Link: "/cp/" + wrap.CurrModule + "/" + wrap.CurrSubModule + "/"},
				{Name: "List of categories"},
			})
			content += builder.DataTable(
				wrap,
				"shop_cats",
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
							Href:   `/shop/category/` + (*values)[3] + `/`,
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
							Href: "javascript:fave.ActionDataTableDelete(this,'shop-categories-delete','" +
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
				},
				false,
			)
		} else if wrap.CurrSubModule == "attributes" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Attributes", Link: "/cp/" + wrap.CurrModule + "/" + wrap.CurrSubModule + "/"},
				{Name: "List of attributes"},
			})
			content += builder.DataTable(
				wrap,
				"shop_filters",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "name",
						NameInTable: "Name",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/attributes-modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
							return `<div>` + name + `</div><div><small>` + html.EscapeString((*values)[2]) + `</small></div>`
						},
					},
					{
						DBField: "filter",
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/attributes-modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'shop-attributes-delete','" +
								(*values)[0] + "','Are you sure want to delete attribute?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/"+wrap.CurrSubModule+"/",
				nil,
				nil,
				true,
			)
		} else if wrap.CurrSubModule == "currencies" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Currencies", Link: "/cp/" + wrap.CurrModule + "/" + wrap.CurrSubModule + "/"},
				{Name: "List of currencies"},
			})
			content += builder.DataTable(
				wrap,
				"shop_currencies",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField:     "name",
						NameInTable: "Name",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/currencies-modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + ` (` + (*values)[3] + `, ` + (*values)[4] + `)</a>`
							return `<div>` + name + `</div>`
						},
					},
					{
						DBField:     "coefficient",
						NameInTable: "Coefficient",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							return utils.Float64ToStrF(utils.StrToFloat64((*values)[2]), "%.4f")
						},
					},
					{
						DBField: "code",
					},
					{
						DBField: "symbol",
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/currencies-modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'shop-currencies-delete','" +
								(*values)[0] + "','Are you sure want to delete currency?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/"+wrap.CurrSubModule+"/",
				nil,
				nil,
				true,
			)
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new product"},
				})
			} else {
				if len(wrap.UrlArgs) >= 3 && utils.IsNumeric(wrap.UrlArgs[2]) {
					content += `<div class="product-copy"><a title="Duplicate product" href="javascript:fave.ShopProductsDuplicate(this, ` + wrap.UrlArgs[2] + `);">` + assets.SysSvgIconCopy + `</a></div>`
					content += `<div class="product-another"><a title="The same with another attributes" href="javascript:fave.ShopProductsAnother(this, ` + wrap.UrlArgs[2] + `);">` + assets.SysSvgIconPlus + `</a></div>`
				}
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify product"},
				})
			}

			data := utils.MySql_shop_product{
				A_id:       0,
				A_user:     0,
				A_currency: 0,
				A_price:    0,
				A_name:     "",
				A_alias:    "",
				A_vendor:   "",
				A_quantity: 0,
				A_category: 0,
				A_briefly:  "",
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
				err := wrap.DB.QueryRow(`
					SELECT
						id,
						user,
						currency,
						price,
						name,
						alias,
						vendor,
						quantity,
						category,
						briefly,
						content,
						active
					FROM
						shop_products
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_currency,
					&data.A_price,
					&data.A_name,
					&data.A_alias,
					&data.A_vendor,
					&data.A_quantity,
					&data.A_category,
					&data.A_briefly,
					&data.A_content,
					&data.A_active,
				)
				if err != nil {
					return "", "", ""
				}
			}

			// All product current categories
			var selids []int
			if data.A_id > 0 {
				rows, err := wrap.DB.Query("SELECT category_id FROM shop_cat_product_rel WHERE product_id = ?;", data.A_id)
				if err == nil {
					defer rows.Close()
					values := make([]int, 1)
					scan := make([]interface{}, len(values))
					for i := range values {
						scan[i] = &values[i]
					}
					for rows.Next() {
						err = rows.Scan(scan...)
						if err == nil {
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
					Value: "shop-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Product name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Product price",
					Name:    "price",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price">Product price</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-8">` +
							`<div><input class="form-control" type="number" step="0.01" id="lbl_price" name="price" value="` + utils.Float64ToStr(data.A_price) + `" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-4">` +
							`<select class="selectpicker form-control" id="lbl_currency" name="currency" data-live-search="true">` +
							this.shop_GetCurrencySelectOptions(wrap, data.A_currency) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKText,
					Caption: "Product alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: mobile-phone",
					Max:     "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Vendor/Count",
					Name:    "vendor",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n5">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_vendor">Vendor/Count</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-8">` +
							`<div><input class="form-control" type="text" id="lbl_vendor" name="vendor" value="` + html.EscapeString(data.A_vendor) + `" placeholder="" autocomplete="off"></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-4">` +
							`<input class="form-control" type="number" step="1" id="lbl_quantity" name="quantity" value="` + utils.IntToStr(data.A_quantity) + `" placeholder="" autocomplete="off">` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKText,
					Caption: "Category",
					Name:    "category",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n6">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_category">Category</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_category" name="category" data-live-search="true">` +
							`<option title="Nothing selected" value="0">&mdash;</option>` +
							this.shop_GetCategorySelectOptions(wrap, 0, data.A_category, []int{}) +
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
						return `<div class="form-group n7">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_cats">Categories</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_cats" name="cats[]" data-live-search="true" multiple>` +
							this.shop_GetCategorySelectOptions(wrap, 0, 0, selids) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKText,
					Caption: "Attributes",
					Name:    "",
					Value:   "",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n8">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Attributes</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div class="list-wrapper">` +
							`<div id="list">` +
							this.shop_GetProductValuesInputs(wrap, data.A_id) +
							`</div>` +
							`<div class="list-button position-relative">` +
							`<select class="selectpicker form-control" id="lbl_attributes" data-live-search="true" onchange="fave.ShopProductsAdd();">` +
							this.shop_GetAllAttributesSelectOptions(wrap) +
							`</select>` +
							`</div>` +
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
					Caption: "Product content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "wysiwyg",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Product images",
					Name:    "",
					Value:   "",
					CallBack: func(field *builder.DataFormField) string {
						if data.A_id == 0 {
							return ``
						}
						return `<div class="form-group n11">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Product images</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div class="list-wrapper">` +
							`<div id="list-images">` +
							this.shop_GetAllProductImages(wrap, data.A_id) +
							`</div>` +
							`<div id="img-upload-block" class="list-button position-relative">` +
							`<div id="upload-msg">Uploading...</div>` +
							`<input class="form-control ignore-lost-data" type="file" id="file" name="file" onchange="fave.ShopProductsUploadImage('shop-upload-image', ` + utils.IntToStr(data.A_id) + `, 'file');" style="font-size:13px;" multiple />` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`<script>WaitForFave(function(){Sortable.create(document.getElementById('list-images'),{animation:0,onEnd:function(evt){var orderData=[];$('#list-images div.attached-img').each(function(i,v){orderData.push({Id:$(v).data('id'),Order:i+1});});$('#list-images').addClass('loading');fave.ShopProductsImageReorder('shop-images-reorder',{Items:orderData});},});});</script>`
					},
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

			data := utils.MySql_shop_category{
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
						shop_cats
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

			btn_caption := "Add"
			if wrap.CurrSubModule == "categories-modify" {
				btn_caption = "Save"
			}

			parentId := 0
			if wrap.CurrSubModule == "categories-modify" {
				parentId = this.shop_GetCategoryParentId(wrap, data.A_id)
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "shop-categories-modify",
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
							this.shop_GetCategorySelectOptions(wrap, data.A_id, parentId, []int{}) +
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
					Hint:    "Example: popular-products",
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
		} else if wrap.CurrSubModule == "attributes-add" || wrap.CurrSubModule == "attributes-modify" {
			if wrap.CurrSubModule == "attributes-add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Attributes", Link: "/cp/" + wrap.CurrModule + "/attributes/"},
					{Name: "Add new attribute"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Attributes", Link: "/cp/" + wrap.CurrModule + "/attributes/"},
					{Name: "Modify attribute"},
				})
			}

			data := utils.MySql_shop_filter{
				A_id:     0,
				A_name:   "",
				A_filter: "",
			}

			if wrap.CurrSubModule == "attributes-modify" {
				if len(wrap.UrlArgs) != 3 {
					return "", "", ""
				}
				if !utils.IsNumeric(wrap.UrlArgs[2]) {
					return "", "", ""
				}
				err := wrap.DB.QueryRow(`
					SELECT
						id,
						name,
						filter
					FROM
						shop_filters
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_name,
					&data.A_filter,
				)
				if err != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "attributes-modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "shop-attributes-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Attribute name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Attribute in filter",
					Name:     "filter",
					Value:    data.A_filter,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Attribute values",
					Name:    "",
					Value:   "",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n4">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Attribute values</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div class="list-wrapper">` +
							`<div id="list">` +
							this.shop_GetFilterValuesInputs(wrap, data.A_id) +
							`</div>` +
							`<div class="list-button"><button type="button" class="btn btn-success" onclick="fave.ShopAttributesAdd();">Add attribute value</button></div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  btn_caption,
					Target: "add-edit-button",
				},
			})

			if wrap.CurrSubModule == "attributes-add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
		} else if wrap.CurrSubModule == "currencies-add" || wrap.CurrSubModule == "currencies-modify" {
			if wrap.CurrSubModule == "currencies-add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Currencies", Link: "/cp/" + wrap.CurrModule + "/currencies/"},
					{Name: "Add new currency"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Currencies", Link: "/cp/" + wrap.CurrModule + "/currencies/"},
					{Name: "Modify currency"},
				})
			}

			data := utils.MySql_shop_currency{
				A_id:          0,
				A_name:        "",
				A_coefficient: 0,
				A_code:        "",
				A_symbol:      "",
			}

			if wrap.CurrSubModule == "currencies-modify" {
				if len(wrap.UrlArgs) != 3 {
					return "", "", ""
				}
				if !utils.IsNumeric(wrap.UrlArgs[2]) {
					return "", "", ""
				}
				err := wrap.DB.QueryRow(`
					SELECT
						id,
						name,
						coefficient,
						code,
						symbol
					FROM
						shop_currencies
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_name,
					&data.A_coefficient,
					&data.A_code,
					&data.A_symbol,
				)
				if err != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "currencies-modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "shop-currencies-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Currency name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Currency coefficient",
					Name:    "coefficient",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_coefficient">Currency coefficient</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div><input class="form-control" type="number" step="0.0001" id="lbl_coefficient" name="coefficient" value="` + utils.Float64ToStrF(data.A_coefficient, "%.4f") + `" placeholder="" autocomplete="off" required></div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Currency code",
					Name:     "code",
					Value:    data.A_code,
					Required: true,
					Min:      "1",
					Max:      "10",
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Currency symbol",
					Name:     "symbol",
					Value:    data.A_symbol,
					Required: true,
					Min:      "1",
					Max:      "5",
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  btn_caption,
					Target: "add-edit-button",
				},
			})

			if wrap.CurrSubModule == "currencies-add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
