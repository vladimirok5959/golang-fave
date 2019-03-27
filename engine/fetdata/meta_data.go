package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) MetaTitle() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_title
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return ""
			} else if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] != "" {
				// Blog post
				return ""
			} else {
				// Blog
				return ""
			}
		}
	}
	return ""
}

func (this *FERData) MetaKeywords() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_keywords
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return ""
			} else if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] != "" {
				// Blog post
				return ""
			} else {
				// Blog
				return ""
			}
		}
	}
	return ""
}

func (this *FERData) MetaDescription() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_meta_description
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return ""
			} else if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] != "" {
				// Blog post
				return ""
			} else {
				// Blog
				return ""
			}
		}
	}
	return ""
}
