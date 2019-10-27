package basket

import (
	"encoding/json"
	"html"
	"strings"

	"golang-fave/engine/sqlw"
	"golang-fave/utils"
)

type session struct {
	listCurrencies map[int]*currency
	totalSum       float64

	Products   map[int]*product `json:"products"`
	Currency   *currency        `json:"currency"`
	TotalSum   string           `json:"total_sum"`
	TotalCount int              `json:"total_count"`
}

func (this *session) makePrice(product_price float64, product_currency_id int) float64 {
	if this.Currency == nil {
		return product_price
	}
	if this.Currency.Id == product_currency_id {
		return product_price
	}
	if product_currency_id == 1 {
		return product_price * this.Currency.Coefficient
	} else {
		if c, ok := this.listCurrencies[product_currency_id]; ok == true {
			return product_price / c.Coefficient
		} else {
			return product_price
		}
	}
}

func (this *session) updateProducts(db *sqlw.DB) {
	products_ids := []int{}
	for _, product := range this.Products {
		products_ids = append(products_ids, product.Id)
	}
	if len(products_ids) > 0 {
		if rows, err := db.Query(
			`SELECT
				shop_products.id,
				shop_products.name,
				shop_products.price,
				shop_products.alias,
				shop_currencies.id,
				shop_currencies.name,
				shop_currencies.coefficient,
				shop_currencies.code,
				shop_currencies.symbol,
				IF(image_this.filename IS NULL, IFNULL(shop_products.parent_id, shop_products.id), shop_products.id) as imgid,
				IFNULL(IFNULL(image_this.filename, image_parent.filename), '') as filename
			FROM
				shop_products
				LEFT JOIN shop_currencies ON shop_currencies.id = shop_products.currency
				LEFT JOIN (
					SELECT
						m.product_id,
						m.filename
					FROM
						shop_product_images as m
						LEFT JOIN (
							SELECT
								t.product_id,
								MIN(t.ord) as ordmin
							FROM
								shop_product_images as t
							GROUP BY
								t.product_id
						) as u ON u.product_id = m.product_id AND u.ordmin = m.ord
					WHERE
						u.product_id IS NOT NULL
				) as image_this ON image_this.product_id = shop_products.id
				LEFT JOIN (
					SELECT
						m.product_id,
						m.filename
					FROM
						shop_product_images as m
						LEFT JOIN (
							SELECT
								t.product_id,
								MIN(t.ord) as ordmin
							FROM
								shop_product_images as t
							GROUP BY
								t.product_id
						) as u ON u.product_id = m.product_id AND u.ordmin = m.ord
					WHERE
						u.product_id IS NOT NULL
				) as image_parent ON image_parent.product_id = shop_products.parent_id
			WHERE
				shop_products.active = 1 AND
				shop_products.id IN (` + strings.Join(utils.ArrayOfIntToArrayOfString(products_ids), ",") + `)
			;`,
		); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := &utils.MySql_shop_product{}
				roc := &utils.MySql_shop_currency{}
				var img_product_id string
				var img_filename string
				if err = rows.Scan(
					&row.A_id,
					&row.A_name,
					&row.A_price,
					&row.A_alias,
					&roc.A_id,
					&roc.A_name,
					&roc.A_coefficient,
					&roc.A_code,
					&roc.A_symbol,
					&img_product_id,
					&img_filename,
				); err == nil {
					if p, ok := this.Products[row.A_id]; ok == true {
						var product_image string
						if img_filename == "" {
							product_image = utils.GetImagePlaceholderSrc()
						} else {
							product_image = "/products/images/" + img_product_id + "/thumb-0-" + img_filename
						}
						p.Name = html.EscapeString(row.A_name)
						p.Image = product_image
						p.Link = "/shop/" + row.A_alias + "/"
						p.price = row.A_price
						p.currency.Id = roc.A_id
						p.currency.Name = html.EscapeString(roc.A_name)
						p.currency.Coefficient = roc.A_coefficient
						p.currency.Code = html.EscapeString(roc.A_code)
						p.currency.Symbol = html.EscapeString(roc.A_symbol)
					}
				}
			}
		}
	}
}

func (this *session) updateTotals(p *SBParam) {
	this.totalSum = 0
	this.TotalCount = 0
	for _, product := range this.Products {
		product.Price = utils.FormatProductPrice(this.makePrice(product.price, product.currency.Id), (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
		product.Sum = utils.FormatProductPrice(this.makePrice(product.price*float64(product.Quantity), product.currency.Id), (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
		this.totalSum += this.makePrice(product.price, product.currency.Id) * float64(product.Quantity)
		this.TotalCount += product.Quantity
	}
	this.TotalSum = utils.FormatProductPrice(this.totalSum, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
}

func (this *session) Preload(p *SBParam) {
	user_currency := 1
	if cookie, err := p.R.Cookie("currency"); err == nil {
		user_currency = utils.StrToInt(cookie.Value)
	}

	// Clear list of currencies
	this.listCurrencies = map[int]*currency{}

	// Load currencies from database
	if rows, err := p.DB.Query(
		`SELECT
			id,
			name,
			coefficient,
			code,
			symbol
		FROM
			shop_currencies
		ORDER BY
			id ASC
		;`,
	); err == nil {
		defer rows.Close()
		for rows.Next() {
			roc := &utils.MySql_shop_currency{}
			if err = rows.Scan(
				&roc.A_id,
				&roc.A_name,
				&roc.A_coefficient,
				&roc.A_code,
				&roc.A_symbol,
			); err == nil {
				this.listCurrencies[roc.A_id] = &currency{
					Id:          roc.A_id,
					Name:        html.EscapeString(roc.A_name),
					Coefficient: roc.A_coefficient,
					Code:        html.EscapeString(roc.A_code),
					Symbol:      html.EscapeString(roc.A_symbol),
				}
			}
		}
	}

	// Check if selected currency is exists
	if _, ok := this.listCurrencies[user_currency]; ok != true {
		user_currency = 1
	}

	// Remember selected currency
	if c, ok := this.listCurrencies[user_currency]; ok == true {
		this.Currency = &currency{
			Id:          c.Id,
			Name:        c.Name,
			Coefficient: c.Coefficient,
			Code:        c.Code,
			Symbol:      c.Symbol,
		}
	}
}

func (this *session) String(p *SBParam) string {
	this.updateProducts(p.DB)
	this.updateTotals(p)

	json, err := json.Marshal(this)
	if err != nil {
		return `{"msg":"basket_engine_error","message":"` + err.Error() + `"}`
	}

	return string(json)
}

func (this *session) Plus(p *SBParam, product_id int) {
	if prod, ok := this.Products[product_id]; ok == true {
		prod.Quantity++
		this.updateProducts(p.DB)
		this.updateTotals(p)
		return
	}
	row := &utils.MySql_shop_product{}
	roc := &utils.MySql_shop_currency{}
	var img_product_id string
	var img_filename string
	if err := p.DB.QueryRow(`
		SELECT
			shop_products.id,
			shop_products.name,
			shop_products.price,
			shop_products.alias,
			shop_currencies.id,
			shop_currencies.name,
			shop_currencies.coefficient,
			shop_currencies.code,
			shop_currencies.symbol,
			IF(image_this.filename IS NULL, IFNULL(shop_products.parent_id, shop_products.id), shop_products.id) as imgid,
			IFNULL(IFNULL(image_this.filename, image_parent.filename), '') as filename
		FROM
			shop_products
			LEFT JOIN shop_currencies ON shop_currencies.id = shop_products.currency
			LEFT JOIN (
				SELECT
					m.product_id,
					m.filename
				FROM
					shop_product_images as m
					LEFT JOIN (
						SELECT
							t.product_id,
							MIN(t.ord) as ordmin
						FROM
							shop_product_images as t
						GROUP BY
							t.product_id
					) as u ON u.product_id = m.product_id AND u.ordmin = m.ord
				WHERE
					u.product_id IS NOT NULL
			) as image_this ON image_this.product_id = shop_products.id
			LEFT JOIN (
				SELECT
					m.product_id,
					m.filename
				FROM
					shop_product_images as m
					LEFT JOIN (
						SELECT
							t.product_id,
							MIN(t.ord) as ordmin
						FROM
							shop_product_images as t
						GROUP BY
							t.product_id
					) as u ON u.product_id = m.product_id AND u.ordmin = m.ord
				WHERE
					u.product_id IS NOT NULL
			) as image_parent ON image_parent.product_id = shop_products.parent_id
		WHERE
			shop_products.active = 1 AND
			shop_products.id = ?
		LIMIT 1;`,
		product_id,
	).Scan(
		&row.A_id,
		&row.A_name,
		&row.A_price,
		&row.A_alias,
		&roc.A_id,
		&roc.A_name,
		&roc.A_coefficient,
		&roc.A_code,
		&roc.A_symbol,
		&img_product_id,
		&img_filename,
	); err == nil {
		var product_image string
		if img_filename == "" {
			product_image = utils.GetImagePlaceholderSrc()
		} else {
			product_image = "/products/images/" + img_product_id + "/thumb-0-" + img_filename
		}
		this.Products[product_id] = &product{
			currency: &currency{Id: roc.A_id, Name: roc.A_name, Coefficient: roc.A_coefficient, Code: roc.A_code, Symbol: roc.A_symbol},

			Id:       row.A_id,
			Name:     html.EscapeString(row.A_name),
			Image:    product_image,
			Link:     "/shop/" + row.A_alias + "/",
			price:    row.A_price,
			Quantity: 1,
		}
		this.updateProducts(p.DB)
		this.updateTotals(p)
	}
}

func (this *session) Minus(p *SBParam, product_id int) {
	if prod, ok := this.Products[product_id]; ok == true {
		if prod.Quantity > 1 {
			prod.Quantity--
		} else {
			delete(this.Products, product_id)
		}
		this.updateProducts(p.DB)
		this.updateTotals(p)
	}
}

func (this *session) Remove(p *SBParam, product_id int) {
	if _, ok := this.Products[product_id]; ok == true {
		delete(this.Products, product_id)
		this.updateProducts(p.DB)
		this.updateTotals(p)
	}
}

func (this *session) ProductsCount() int {
	return this.TotalCount
}
