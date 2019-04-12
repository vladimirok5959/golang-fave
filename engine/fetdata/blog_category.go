package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type BlogCategory struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_blog_category
	depth  int

	user *User
}

func (this *BlogCategory) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *BlogCategory) User() *User {
	if this == nil {
		return nil
	}
	if this.user != nil {
		return this.user
	}
	this.user = &User{wrap: this.wrap}
	this.user.load(this.object.A_user)
	return this.user
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

func (this *BlogCategory) Level() int {
	if this == nil {
		return 0
	}
	return this.depth
}
