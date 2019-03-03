package fetdata

import (
	"html/template"
	"time"

	"golang-fave/utils"
)

func (this *FERData) Name() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_name
		}
	}
	return ""
}

func (this *FERData) Alias() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_alias
		}
	}
	return ""
}

func (this *FERData) Content() template.HTML {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return template.HTML(this.dataRow.(*utils.MySql_page).A_content)
		}
	}
	return template.HTML("")
}

func (this *FERData) DateTime() int {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_datetime
		}
	}
	return 0
}

func (this *FERData) DateTimeFormat(format string) string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return time.Unix(int64(this.dataRow.(*utils.MySql_page).A_datetime), 0).Format(format)
		}
	}
	return ""
}
