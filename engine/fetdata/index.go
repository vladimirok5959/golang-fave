package fetdata

import (
	"html/template"
	"time"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

type Page struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_page

	user *User
}

func (this *Page) load() *Page {
	return this
}

func (this *Page) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *Page) User() *User {
	if this == nil {
		return nil
	}
	if this.user != nil {
		return this.user
	}
	this.user = &User{wrap: this.wrap}
	this.user.loadById(this.object.A_user)
	return this.user
}

func (this *Page) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *Page) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *Page) Content() template.HTML {
	if this == nil {
		return template.HTML("")
	}
	return template.HTML(this.object.A_content)
}

func (this *Page) MetaTitle() string {
	if this == nil {
		return ""
	}
	return this.object.A_meta_title
}

func (this *Page) MetaKeywords() string {
	if this == nil {
		return ""
	}
	return this.object.A_meta_keywords
}

func (this *Page) MetaDescription() string {
	if this == nil {
		return ""
	}
	return this.object.A_meta_description
}

func (this *Page) DateTimeUnix() int {
	if this == nil {
		return 0
	}
	return this.object.A_datetime
}

func (this *Page) DateTimeFormat(format string) string {
	if this == nil {
		return ""
	}
	return time.Unix(int64(this.object.A_datetime), 0).Format(format)
}

func (this *Page) Active() bool {
	if this == nil {
		return false
	}
	return this.object.A_active > 0
}
