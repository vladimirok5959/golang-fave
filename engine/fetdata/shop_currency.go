package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopCurrency struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_currency
}

func (this *ShopCurrency) load() *ShopCurrency {
	return this
}

func (this *ShopCurrency) loadById(id int) {
	if this == nil {
		return
	}
	if this.object != nil {
		return
	}
	this.object = &utils.MySql_shop_currency{}
	if err := this.wrap.DB.QueryRow(
		this.wrap.R.Context(),
		`SELECT
			id,
			name,
			coefficient,
			code,
			symbol
		FROM
			fave_shop_currencies
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
	); *this.wrap.LogCpError(&err) != nil {
		return
	}
}

func (this *ShopCurrency) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *ShopCurrency) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *ShopCurrency) Coefficient() float64 {
	if this == nil {
		return 0
	}
	return this.object.A_coefficient
}

func (this *ShopCurrency) Code() string {
	if this == nil {
		return ""
	}
	return this.object.A_code
}

func (this *ShopCurrency) Symbol() string {
	if this == nil {
		return ""
	}
	return this.object.A_symbol
}
