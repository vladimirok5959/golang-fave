package fetdata

import (
	"html/template"
	"time"

	"golang-fave/utils"
)

func (this *FERData) Name() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_name
		}
	}
	return ""
}

func (this *FERData) Alias() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_alias
		}
	}
	return ""
}

func (this *FERData) Content() template.HTML {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return template.HTML(this.DataRow.(*utils.MySql_page).A_content)
		}
	}
	return template.HTML("")
}

func (this *FERData) DateTime() int {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_datetime
		}
	}
	return 0
}

func (this *FERData) DateTimeFormat(format string) string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return time.Unix(int64(this.DataRow.(*utils.MySql_page).A_datetime), 0).Format(format)
		}
	}
	return ""
}
