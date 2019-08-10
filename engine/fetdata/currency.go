package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type Currency struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_currency
}

func (this *Currency) load(id int) {
	if this == nil {
		return
	}
	if this.object != nil {
		return
	}
	this.object = &utils.MySql_shop_currency{}
	if err := this.wrap.DB.QueryRow(`
		SELECT
			id,
			name,
			coefficient,
			code,
			symbol
		FROM
			shop_currencies
		WHERE
			id = ?
		LIMIT 1;`,
		id,
	).Scan(
		&this.object.A_id,
		&this.object.A_name,
		&this.object.A_coefficient,
		&this.object.A_code,
		&this.object.A_symbol,
	); err != nil {
		return
	}
}

func (this *Currency) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *Currency) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *Currency) Coefficient() float64 {
	if this == nil {
		return 0
	}
	return this.object.A_coefficient
}

func (this *Currency) Code() string {
	if this == nil {
		return ""
	}
	return this.object.A_code
}

func (this *Currency) Symbol() string {
	if this == nil {
		return ""
	}
	return this.object.A_symbol
}
