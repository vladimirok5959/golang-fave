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

func (this *ShopCategory) load() *ShopCategory {
	return this
}

func (this *ShopCategory) loadById(id int) {
	if this == nil {
		return
	}
	if this.object != nil {
		return
	}
	this.object = &utils.MySql_shop_category{}
	if err := this.wrap.DB.QueryRow(`
		SELECT
			id,
			user,
			name,
			alias,
			lft,
			rgt
		FROM
			shop_cats
		WHERE
			id = ? AND
			id > 1
		LIMIT 1;`,
		id,
	).Scan(
		&this.object.A_id,
		&this.object.A_user,
		&this.object.A_name,
		&this.object.A_alias,
		&this.object.A_lft,
		&this.object.A_rgt,
	); err != nil {
		return
	}
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
	this.user.loadById(this.object.A_user)
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
