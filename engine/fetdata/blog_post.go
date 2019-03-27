package fetdata

import (
	"html/template"
	"time"
)

type BlogPost struct {
	id       int
	user     int
	name     string
	alias    string
	content  string
	datetime int
	active   int
}

func (this *BlogPost) Id() int {
	return this.id
}

func (this *BlogPost) Name() string {
	return this.name
}

func (this *BlogPost) Alias() string {
	return this.alias
}

func (this *BlogPost) Permalink() string {
	return "/blog/" + this.alias + "/"
}

func (this *BlogPost) Content() template.HTML {
	return template.HTML(this.content)
}

func (this *BlogPost) DateTime() int {
	return this.datetime
}

func (this *BlogPost) DateTimeFormat(format string) string {
	return time.Unix(int64(this.datetime), 0).Format(format)
}
