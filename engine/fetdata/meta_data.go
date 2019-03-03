package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) MetaTitle() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_title
		}
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_keywords
		}
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_description
		}
	}
	return ""
}
