package basket

import (
	"encoding/json"
	"html"

	"golang-fave/engine/sqlw"
	"golang-fave/utils"
)

type session struct {
	Products   map[int]*product `json:"products"`
	Currency   *currency        `json:"currency"`
	TotalSum   float64          `json:"total_sum"`
	TotalCount int              `json:"total_count"`
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
		p.Sum = p.Price * float64(p.Quantity)
		this.updateTotals()
		return
	}
	row := &utils.MySql_shop_product{}
	if err := db.QueryRow(`
		SELECT
			shop_products.id,
			shop_products.name,
			shop_products.price
		FROM
			shop_products
		WHERE
			shop_products.active = 1 AND
			shop_products.id = ?
		LIMIT 1;`,
		product_id,
	).Scan(
		&row.A_id,
		&row.A_name,
		&row.A_price,
	); err == nil {
		// Load product image here
		this.Products[product_id] = &product{
			Id:       row.A_id,
			Name:     html.EscapeString(row.A_name),
			Image:    "",
			Price:    row.A_price,
			Quantity: 1,
			Sum:      row.A_price,
		}
		this.updateTotals()
	}
}

func (this *session) Minus(db *sqlw.DB, product_id int) {
	if p, ok := this.Products[product_id]; ok == true {
		if p.Quantity > 1 {
			p.Quantity--
			p.Sum = p.Price * float64(p.Quantity)
		} else {
			delete(this.Products, product_id)
		}
		this.updateTotals()
	}
}

func (this *session) Remove(db *sqlw.DB, product_id int) {
	if _, ok := this.Products[product_id]; ok == true {
		delete(this.Products, product_id)
		this.updateTotals()
	}
}

func (this *session) updateTotals() {
	this.TotalSum = 0
	this.TotalCount = 0
	for _, product := range this.Products {
		this.TotalSum += product.Price * float64(product.Quantity)
		this.TotalCount += product.Quantity
	}
}
