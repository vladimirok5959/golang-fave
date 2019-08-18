package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopCategory struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_category

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
			main.id,
			main.user,
			main.name,
			main.alias,
			main.lft,
			main.rgt,
			depth.depth,
			MAX(main.parent_id) AS parent_id
		FROM
			(
				SELECT
					node.id,
					node.user,
					node.name,
					node.alias,
					node.lft,
					node.rgt,
					parent.id AS parent_id
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt AND
					node.id > 1
				ORDER BY
					node.lft ASC
			) AS main
			LEFT JOIN (
				SELECT
					node.id,
					(COUNT(parent.id) - 1) AS depth
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt
				GROUP BY
					node.id
				ORDER BY
					node.lft ASC
			) AS depth ON depth.id = main.id
		WHERE
			main.id > 1 AND
			main.id <> main.parent_id AND
			main.id = ?
		GROUP BY
			main.id
		;`,
		id,
	).Scan(
		&this.object.A_id,
		&this.object.A_user,
		&this.object.A_name,
		&this.object.A_alias,
		&this.object.A_lft,
		&this.object.A_rgt,
		&this.object.A_depth,
		&this.object.A_parent,
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
	this.user = (&User{wrap: this.wrap}).load()
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
	return this.object.A_depth
}

func (this *ShopCategory) ParentId() int {
	if this == nil {
		return 0
	}
	return this.object.A_parent
}
