package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) MetaTitle() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_title
		}
	} else if this.is404 {
		// Return it from settings
		return "Page not found"
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_keywords
		}
	} else if this.is404 {
		// Return it from settings
		return ""
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_description
		}
	} else if this.is404 {
		// Return it from settings
		return ""
	}
	return ""
}
