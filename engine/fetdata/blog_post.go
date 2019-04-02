package fetdata

import (
	"html/template"
	"time"

	"golang-fave/utils"
)

type BlogPost struct {
	object *utils.MySql_blog_posts
}

func (this *BlogPost) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *BlogPost) User() int {
	if this == nil {
		return 0
	}
	return this.object.A_user
}

func (this *BlogPost) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *BlogPost) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *BlogPost) Content() template.HTML {
	if this == nil {
		return template.HTML("")
	}
	return template.HTML(this.object.A_content)
}

func (this *BlogPost) DateTimeUnix() int {
	if this == nil {
		return 0
	}
	return this.object.A_datetime
}

func (this *BlogPost) DateTimeFormat(format string) string {
	if this == nil {
		return ""
	}
	return time.Unix(int64(this.object.A_datetime), 0).Format(format)
}

func (this *BlogPost) Active() bool {
	if this == nil {
		return false
	}
	return this.object.A_active > 0
}

func (this *BlogPost) Permalink() string {
	if this == nil {
		return ""
	}
	return "/blog/" + this.object.A_alias + "/"
}
