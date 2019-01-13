package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
)

type dataTableDisplay func(values *[]string) string

type dataTableAction func(values *[]string) string

type dataTableRow struct {
	dbField     string
	nameInTable string
	display     dataTableDisplay
}

type dataBreadcrumb struct {
	name string
	link string
}

type ModuleItem struct {
	Alias   string
	Display bool
	Name    string
	Icon    string
	Order   int
}

type Module struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
	user    *utils.MySql_user
	urls    *[]string
	mmod    string
	smod    string
	imod    int
	modlist []ModuleItem
}

func (this *Module) module_get_display(name string) bool {
	mname := "Module_" + name + "_display"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].Bool()
	}
	return false
}

func (this *Module) module_get_name(name string) string {
	mname := "Module_" + name + "_name"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].String()
	}
	return ""
}

func (this *Module) module_get_icon(name string) string {
	mname := "Module_" + name + "_icon"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].String()
	}
	return ""
}

func (this *Module) module_get_order(name string) int {
	mname := "Module_" + name + "_order"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return int(result[0].Int())
	}
	return 0
}

func (this *Module) module_get_submenu(name string) string {
	mname := "Module_" + name + "_submenu"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		result_array := result[0].Interface().([]utils.ModuleSubMenu)
		result_html := ""
		for _, value := range result_array {
			class := ""
			if name == this.mmod && value.Alias == this.smod {
				class = " active"
			}
			result_html += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + name + `/` + value.Alias + `/">` + value.Icon + value.Name + `</a></li>`
		}
		if result_html != "" {
			result_html = `<ul class="nav flex-column">` + result_html + `</ul>`
		}
		return result_html
	}
	return ""
}

func (this *Module) module_get_list_of_modules() *[]ModuleItem {
	if len(this.modlist) <= 0 {
		t := reflect.TypeOf(this)
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if strings.HasPrefix(m.Name, "Module_") && strings.HasSuffix(m.Name, "_alias") {
				alias := m.Name[7:]
				alias = alias[0 : len(alias)-6]
				this.modlist = append(this.modlist, ModuleItem{
					alias,
					this.module_get_display(alias),
					this.module_get_name(alias),
					this.module_get_icon(alias),
					this.module_get_order(alias),
				})
			}
		}
		sort.Slice(this.modlist, func(i, j int) bool {
			return this.modlist[i].Order < this.modlist[j].Order
		})
	}
	return &this.modlist
}

func (this *Module) breadcrumb(data []dataBreadcrumb) string {
	result := ``
	result += `<nav aria-label="breadcrumb">`
	result += `<ol class="breadcrumb">`
	result += `<li class="breadcrumb-item"><a href="/cp/` + this.mmod + `/">` + this.module_get_name(this.mmod) + `</a></li>`
	for _, item := range data {
		if item.link == "" {
			result += `<li class="breadcrumb-item active" aria-current="page">` + item.name + `</li>`
		} else {
			result += `<li class="breadcrumb-item"><a href="` + item.link + `">` + item.name + `</a></li>`
		}
	}
	result += `</ol>`
	result += `</nav>`
	return result
}

func (this *Module) data_table(table string, order_by string, order_way string, data []dataTableRow, action dataTableAction, pagination_url string) string {
	var num int
	err := this.db.QueryRow("SELECT COUNT(*) FROM `" + table + "`;").Scan(&num)
	if err != nil {
		return ""
	}
	pear_page := 10
	max_pages := int(math.Ceil(float64(num) / float64(pear_page)))
	curr_page := 1
	p := this.wrapper.R.URL.Query().Get("p")
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
	result := `<table class="table table-striped table-bordered table_` + table + `">`
	result += `<thead>`
	result += `<tr>`
	sqld := "SELECT"
	for i, column := range data {
		if column.nameInTable != "" {
			result += `<th scope="col" class="col_` + column.dbField + `">` + column.nameInTable + `</th>`
		}
		sqld += " `" + column.dbField + "`"
		if i+1 < len(data) {
			sqld += ","
		}
	}
	sqld += " FROM `" + table + "` ORDER BY `" + order_by + "` " + order_way + " LIMIT ?, ?;"
	if action != nil {
		result += `<th scope="col" class="col_action">Action</th>`
	}
	result += `</tr>`
	result += `</thead>`
	result += `<tbody>`
	rows, err := this.db.Query(sqld, limit_offset, pear_page)
	if err == nil {
		values := make([]string, len(data))
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<tr>`
				for i, val := range values {
					if data[i].nameInTable != "" {
						if data[i].display == nil {
							result += `<td class="col_` + data[i].dbField + `">` + string(val) + `</td>`
						} else {
							result += `<td class="col_` + data[i].dbField + `">` + data[i].display(&values) + `</td>`
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
	result += `</tbody></table>`
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
	return result
}

func New(wrapper *wrapper.Wrapper, db *sql.DB, user *utils.MySql_user, url_args *[]string) *Module {
	mmod := "index"
	smod := "default"
	imod := 0
	if len(*url_args) >= 2 {
		mmod = (*url_args)[1]
	}
	if len(*url_args) >= 3 {
		smod = (*url_args)[2]
	}
	if len(*url_args) >= 4 {
		if val, err := strconv.Atoi((*url_args)[3]); err == nil {
			imod = val
		}
	}
	return &Module{wrapper, db, user, url_args, mmod, smod, imod, make([]ModuleItem, 0)}
}

func (this *Module) Run() bool {
	mname := "Module_" + this.mmod
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return true
	}
	return false
}

func (this *Module) GetNavMenuModules() string {
	html := ""
	list := this.module_get_list_of_modules()
	for _, value := range *list {
		if value.Display {
			class := ""
			if value.Alias == this.mmod {
				class = " active"
			}
			html += `<a class="dropdown-item` + class + `" href="/cp/` + value.Alias + `/">` + value.Name + `</a>`
		}
	}
	return html
}

func (this *Module) GetNavMenuModulesSys() string {
	html := ""
	list := this.module_get_list_of_modules()
	for _, value := range *list {
		if !value.Display {
			class := ""
			if value.Alias == this.mmod {
				class = " active"
			}
			html += `<a class="dropdown-item` + class + `" href="/cp/` + value.Alias + `/">` + value.Name + `</a>`
		}
	}
	return html
}

func (this *Module) GetSidebarLeft() string {
	list := this.module_get_list_of_modules()

	modules_all := `<ul class="nav flex-column">`
	for _, value := range *list {
		if value.Display {
			class := ""
			submenu := ""
			if value.Alias == this.mmod {
				class = " active"
				submenu = this.module_get_submenu(value.Alias)
			}
			modules_all += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + value.Alias + `/">` + value.Icon + value.Name + `</a>` + submenu + `</li>`
		}
	}
	modules_all += `</ul>`

	modules_sys := `<ul class="nav flex-column">`
	for _, value := range *list {
		if !value.Display {
			class := ""
			submenu := ""
			if value.Alias == this.mmod {
				class = " active"
				submenu = this.module_get_submenu(value.Alias)
			}
			modules_sys += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + value.Alias + `/">` + value.Icon + value.Name + `</a>` + submenu + `</li>`
		}
	}
	modules_sys += `</ul>`

	return modules_all + `<div class="dropdown-divider"></div>` + modules_sys
}

func (this *Module) GetContent() string {
	mname := "Module_" + this.mmod + "_content"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].String()
	}
	return ""
}

func (this *Module) GetSidebarRight() string {
	mname := "Module_" + this.mmod + "_sidebar"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].String()
	}
	return ""
}
