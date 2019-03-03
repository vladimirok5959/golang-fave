package fetdata

import (
	"html/template"

	"golang-fave/utils"
)

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

func (this *FERData) DateTime() int {
	if this.Wrap.CurrModule == "index" {
		return this.DataRow.(*utils.MySql_page).A_datetime
	}
	return 0
}
