package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) MetaTitle() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_title
		}
	} else if this.Is404 {
		// Return it from settings
		return "Page not found"
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_keywords
		}
	} else if this.Is404 {
		// Return it from settings
		return ""
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_description
		}
	} else if this.Is404 {
		// Return it from settings
		return ""
	}
	return ""
}
