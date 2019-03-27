package fetdata

import (
	"html/template"
	"time"

	"golang-fave/utils"
)

func (this *FERData) Id() int {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_id
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return this.dataRow.(*utils.MySql_blog_category).A_id
			} else {
				// Blog post
				return this.dataRow.(*utils.MySql_blog_posts).A_id
			}
		}
	}
	return 0
}

func (this *FERData) Name() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_name
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return this.dataRow.(*utils.MySql_blog_category).A_name
			} else {
				// Blog post
				return this.dataRow.(*utils.MySql_blog_posts).A_name
			}
		}
	}
	return ""
}

func (this *FERData) Alias() string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_alias
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return this.dataRow.(*utils.MySql_blog_category).A_alias
			} else {
				// Blog post
				return this.dataRow.(*utils.MySql_blog_posts).A_alias
			}
		}
	}
	return ""
}

func (this *FERData) Content() template.HTML {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return template.HTML(this.dataRow.(*utils.MySql_page).A_content)
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return template.HTML("")
			} else {
				// Blog post
				return template.HTML(this.dataRow.(*utils.MySql_blog_posts).A_content)
			}
		}
	}
	return template.HTML("")
}

func (this *FERData) DateTime() int {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return this.dataRow.(*utils.MySql_page).A_datetime
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return 0
			} else {
				// Blog post
				return this.dataRow.(*utils.MySql_blog_posts).A_datetime
			}
		}
	}
	return 0
}

func (this *FERData) DateTimeFormat(format string) string {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			return time.Unix(int64(this.dataRow.(*utils.MySql_page).A_datetime), 0).Format(format)
		} else if this.wrap.CurrModule == "blog" {
			if len(this.wrap.UrlArgs) >= 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" {
				// Blog category
				return ""
			} else {
				// Blog post
				return time.Unix(int64(this.dataRow.(*utils.MySql_blog_posts).A_datetime), 0).Format(format)
			}
		}
	}
	return ""
}
