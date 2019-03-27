package fetdata

import (
	"time"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	wrap    *wrapper.Wrapper
	dataRow interface{}
	is404   bool

	bufferUser  *utils.MySql_user
	bufferPosts map[string][]*BlogPost
}

func New(wrap *wrapper.Wrapper, drow interface{}, is404 bool) *FERData {
	fer := &FERData{
		wrap:    wrap,
		dataRow: drow,
		is404:   is404,
	}
	return fer.init()
}

func (this *FERData) init() *FERData {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			if this.dataRow.(*utils.MySql_page).A_meta_title == "" {
				this.dataRow.(*utils.MySql_page).A_meta_title = this.dataRow.(*utils.MySql_page).A_name
			}
		}
	}
	return this
}

func (this *FERData) Module() string {
	if this.is404 {
		return "404"
	}

	var mod string
	if this.wrap.CurrModule == "index" {
		mod = "index"
	} else if this.wrap.CurrModule == "blog" {
		if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
			mod = "blog-category"
		} else if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] != "" {
			mod = "blog-post"
		} else {
			mod = "blog"
		}
	}

	return mod
}

func (this *FERData) CurrentDateTime() int {
	return int(time.Now().Unix())
}

func (this *FERData) CurrentDateTimeFormat(format string) string {
	return time.Unix(int64(time.Now().Unix()), 0).Format(format)
}
