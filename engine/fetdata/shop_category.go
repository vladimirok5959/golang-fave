package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopCategory struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_category
	depth  int

	user *User
}

func (this *ShopCategory) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *ShopCategory) User() *User {
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

func (this *ShopCategory) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *ShopCategory) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *ShopCategory) Left() int {
	if this == nil {
		return 0
	}
	return this.object.A_lft
}

func (this *ShopCategory) Right() int {
	if this == nil {
		return 0
	}
	return this.object.A_rgt
}

func (this *ShopCategory) Permalink() string {
	if this == nil {
		return ""
	}
	return "/shop/category/" + this.object.A_alias + "/"
}

func (this *ShopCategory) Level() int {
	if this == nil {
		return 0
	}
	return this.depth
}
