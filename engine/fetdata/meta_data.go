package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) MetaTitle() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_title
		}
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_keywords
		}
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			return this.DataRow.(*utils.MySql_page).A_meta_description
		}
	}
	return ""
}
