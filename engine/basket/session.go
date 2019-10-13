package basket

import (
	"encoding/json"

	"golang-fave/engine/sqlw"
)

type session struct {
	Products map[int]*product `json:"products"`
	Currency *currency        `json:"currency"`
	Total    float64          `json:"total"`
}

func (this *session) String() string {
	json, err := json.Marshal(this)
	if err != nil {
		return `{"msg":"basket_engine_error","message":"` + err.Error() + `"}`
	}
	return string(json)
}

func (this *session) Info(db *sqlw.DB, currency_id int) {
	// Update prices
	// Update total
}

func (this *session) Plus(db *sqlw.DB, product_id int) {
	if p, ok := this.Products[product_id]; ok == true {
		p.Quantity++
		return
	}
	// Check and insert
	// Load from DB
}

func (this *session) Minus(db *sqlw.DB, product_id int) {
	if p, ok := this.Products[product_id]; ok == true {
		if p.Quantity > 1 {
			p.Quantity--
		} else {
			delete(this.Products, product_id)
		}
	}
}

func (this *session) Remove(db *sqlw.DB, product_id int) {
	if _, ok := this.Products[product_id]; ok == true {
		delete(this.Products, product_id)
	}
}
