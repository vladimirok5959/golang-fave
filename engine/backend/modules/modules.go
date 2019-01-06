package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"reflect"
	"strconv"
	"strings"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
)

type Module struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
	user    *utils.MySql_user
	urls    *[]string
	mmod    string
	smod    string
	imod    int
}

func (this *Module) module_get_name(name string) string {
	mname := "Module_" + name + "_name"
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		result := reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return result[0].String()
	}
	return ""
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
	return &Module{wrapper, db, user, url_args, mmod, smod, imod}
}

func (this *Module) Run() bool {
	mname := "Module_" + this.mmod
	if _, ok := reflect.TypeOf(this).MethodByName(mname); ok {
		reflect.ValueOf(this).MethodByName(mname).Call([]reflect.Value{})
		return true
	}
	return false
}

func (this *Module) GetSidebarLeft() string {
	sidebar := `<ul class="nav flex-column">`

	// Make module list
	aType := reflect.TypeOf(this)
	for i := 0; i < aType.NumMethod(); i++ {
		aMethod := aType.Method(i)
		if strings.HasPrefix(aMethod.Name, "Module_") && strings.HasSuffix(aMethod.Name, "_alias") {
			// Extract module alias
			alias := aMethod.Name[7:][:5]

			// Item active class
			class := ""
			if alias == this.mmod {
				class = " active"
			}

			// Active item sub menu
			submenu := ""
			if alias == this.mmod {
				submenu = this.module_get_submenu(alias)
			}

			// Add module to list
			sidebar += `<li class="nav-item` + class + `"><a class="nav-link" href="/cp/` + alias + `/">` + this.module_get_name(alias) + `</a>` + submenu + `</li>`
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
