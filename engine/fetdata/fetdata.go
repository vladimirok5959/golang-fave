package fetdata

import (
	"html/template"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	Wrap    *wrapper.Wrapper
	DataRow interface{}
}

func New(wrap *wrapper.Wrapper, drow interface{}) *FERData {
	fer := &FERData{
		Wrap:    wrap,
		DataRow: drow,
	}
	return fer.init()
}

func (this *FERData) init() *FERData {
	if this.Wrap.CurrModule == "index" {
		if this.DataRow.(*utils.MySql_page).A_meta_title == "" {
			this.DataRow.(*utils.MySql_page).A_meta_title = this.DataRow.(*utils.MySql_page).A_name
		}
	}
	return this
}

func (this *FERData) MetaTitle() string {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_meta_title
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_meta_keywords
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_meta_description
	}
	return ""
}

func (this *FERData) Name() string {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_name
	}
	return ""
}

func (this *FERData) Alias() string {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_alias
	}
	return ""
}

func (this *FERData) Content() template.HTML {
	if this.Wrap.CurrModule == "index" {
		return template.HTML(this.DataRow.(*utils.MySql_page).A_content)
	}
	return template.HTML("")
}
