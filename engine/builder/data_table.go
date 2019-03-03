package builder

import (
	"fmt"
	"html"
	"math"
	"strconv"

	"golang-fave/engine/wrapper"
)

type DataTableRow struct {
	DBField     string
	DBExp       string
	NameInTable string
	Classes     string
	CallBack    func(values *[]string) string
}

func DataTable(wrap *wrapper.Wrapper, table string, order_by string, order_way string, data *[]DataTableRow, action func(values *[]string) string, pagination_url string) string {
	var num int
	err := wrap.DB.QueryRow("SELECT COUNT(*) FROM `" + table + "`;").Scan(&num)
	if err != nil {
		return ""
	}
	pear_page := 10
	max_pages := int(math.Ceil(float64(num) / float64(pear_page)))
	curr_page := 1
	p := wrap.R.URL.Query().Get("p")
	if p != "" {
		pi, err := strconv.Atoi(p)
		if err != nil {
			curr_page = 1
		} else {
			if pi < 1 {
				curr_page = 1
			} else if pi > max_pages {
				curr_page = max_pages
			} else {
				curr_page = pi
			}
		}
	}
	limit_offset := curr_page*pear_page - pear_page
	result := `<table id="cp-table-` + table + `" class="table data-table table-striped table-bordered table-hover table_` + table + `">`
	result += `<thead>`
	result += `<tr>`
	sql := "SELECT"
	for i, column := range *data {
		if column.NameInTable != "" {
			result += `<th scope="col" class="col_` + column.DBField + ` ` + column.Classes + `">` + html.EscapeString(column.NameInTable) + `</th>`
		}
		if column.DBExp == "" {
			sql += " `" + column.DBField + "`"
		} else {
			sql += " " + column.DBExp + " as `" + column.DBField + "`"
		}
		if i+1 < len(*data) {
			sql += ","
		}
	}
	sql += " FROM `" + table + "` ORDER BY `" + order_by + "` " + order_way + " LIMIT ?, ?;"
	if action != nil {
		result += `<th scope="col" class="col_action">&nbsp;</th>`
	}
	result += `</tr>`
	result += `</thead>`
	result += `<tbody>`
	if num > 0 {
		rows, err := wrap.DB.Query(sql, limit_offset, pear_page)
		if err == nil {
			values := make([]string, len(*data))
			scan := make([]interface{}, len(values))
			for i := range values {
				scan[i] = &values[i]
			}
			for rows.Next() {
				err = rows.Scan(scan...)
				if err == nil {
					result += `<tr>`
					for i, val := range values {
						if (*data)[i].NameInTable != "" {
							if (*data)[i].CallBack == nil {
								result += `<td class="col_` + (*data)[i].DBField + ` ` + (*data)[i].Classes + `">` + html.EscapeString(string(val)) + `</td>`
							} else {
								result += `<td class="col_` + (*data)[i].DBField + ` ` + (*data)[i].Classes + `">` + (*data)[i].CallBack(&values) + `</td>`
							}
						}
					}
					if action != nil {
						result += `<td class="col_action">` + action(&values) + `</td>`
					}
					result += `</tr>`
				}
			}
		}
	} else {
		result += `<tr><td colspan="50">No any data found</td></tr>`
	}
	result += `</tbody></table>`

	// Show page navigation only if pages more then one
	if max_pages > 1 {
		result += `<nav>`
		result += `<ul class="pagination" style="margin-bottom:0px;">`
		class := ""
		if curr_page <= 1 {
			class = " disabled"
		}
		result += `<li class="page-item` + class + `">`
		result += `<a class="page-link" href="` + pagination_url + `?p=` + fmt.Sprintf("%d", curr_page-1) + `" aria-label="Previous">`
		result += `<span aria-hidden="true">&laquo;</span>`
		result += `<span class="sr-only">Previous</span>`
		result += `</a>`
		result += `</li>`
		for i := 1; i <= max_pages; i++ {
			class = ""
			if i == curr_page {
				class = " active"
			}
			result += `<li class="page-item` + class + `">`
			result += `<a class="page-link" href="` + pagination_url + `?p=` + fmt.Sprintf("%d", i) + `">` + fmt.Sprintf("%d", i) + `</a>`
			result += `</li>`
		}
		class = ""
		if curr_page >= max_pages {
			class = " disabled"
		}
		result += `<li class="page-item` + class + `">`
		result += `<a class="page-link" href="` + pagination_url + `?p=` + fmt.Sprintf("%d", curr_page+1) + `" aria-label="Next">`
		result += `<span aria-hidden="true">&raquo;</span>`
		result += `<span class="sr-only">Next</span>`
		result += `</a>`
		result += `</li>`
		result += `</ul>`
		result += `</nav>`
	}

	return result
}
