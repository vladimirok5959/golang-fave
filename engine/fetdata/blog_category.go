package fetdata

import (
	"golang-fave/utils"
)

type BlogCategory struct {
	object *utils.MySql_blog_category
}

func (this *BlogCategory) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *BlogCategory) User() int {
	if this == nil {
		return 0
	}
	return this.object.A_user
}

func (this *BlogCategory) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *BlogCategory) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *BlogCategory) Left() int {
	if this == nil {
		return 0
	}
	return this.object.A_lft
}

func (this *BlogCategory) Right() int {
	if this == nil {
		return 0
	}
	return this.object.A_rgt
}

func (this *BlogCategory) Permalink() string {
	if this == nil {
		return ""
	}
	return "/blog/category/" + this.object.A_alias + "/"
}
