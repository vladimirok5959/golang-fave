package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
)

type ModuleItem struct {
	Alias   string
	Display bool
	Name    string
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
			result_html += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + name + `/` + value.Alias + `/">` + value.Name + `</a></li>`
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
				alias := m.Name[7:][:5]
				this.modlist = append(this.modlist, ModuleItem{
					alias,
					this.module_get_display(alias),
					this.module_get_name(alias),
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

func (this *Module) GetSidebarLeft() string {
	sidebar := `<ul class="nav flex-column">`
	list := this.module_get_list_of_modules()
	for _, value := range *list {
		if value.Display {
			class := ""
			submenu := ""
			if value.Alias == this.mmod {
				class = " active"
				submenu = this.module_get_submenu(value.Alias)
			}
			sidebar += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + value.Alias + `/">` + value.Name + `</a>` + submenu + `</li>`
		}
	}
	sidebar += `</ul>`
	return sidebar
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
