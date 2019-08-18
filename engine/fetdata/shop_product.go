package fetdata

import (
	"html/template"
	"strings"
	"time"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopProduct struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_product

	user     *User
	currency *Currency
	category *ShopCategory

	images []*ShopProductImage
	specs  []*ShopProductSpec
}

func (this *ShopProduct) load() *ShopProduct {
	if this == nil {
		return this
	}
	if rows, err := this.wrap.DB.Query(
		`SELECT
			shop_product_images.product_id,
			shop_product_images.filename
		FROM
			shop_product_images
		WHERE
			shop_product_images.product_id = ?
		ORDER BY
			shop_product_images.filename ASC
		;`,
		this.object.A_id,
	); err == nil {
		defer rows.Close()
		for rows.Next() {
			img := utils.MySql_shop_product_image{}
			if err := rows.Scan(
				&img.A_product_id,
				&img.A_filename,
			); err == nil {
				this.images = append(this.images, &ShopProductImage{wrap: this.wrap, object: &img})
			}
		}
	}

	filter_ids := []int{}
	filter_names := map[int]string{}
	filter_values := map[int][]string{}
	if rows, err := this.wrap.DB.Query(
		`SELECT
			shop_filters.id,
			shop_filters.filter,
			shop_filters_values.name
		FROM
			shop_filter_product_values
			LEFT JOIN shop_filters_values ON shop_filters_values.id = shop_filter_product_values.filter_value_id
			LEFT JOIN shop_filters ON shop_filters.id = shop_filters_values.filter_id
		WHERE
			shop_filter_product_values.product_id = ?
		ORDER BY
			shop_filters.filter ASC,
			shop_filters_values.name ASC
		;`,
		this.object.A_id,
	); err == nil {
		defer rows.Close()
		values := make([]string, 3)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				if !utils.InArrayInt(filter_ids, utils.StrToInt(string(values[0]))) {
					filter_ids = append(filter_ids, utils.StrToInt(string(values[0])))
				}
				filter_names[utils.StrToInt(string(values[0]))] = string(values[1])
				filter_values[utils.StrToInt(string(values[0]))] = append(filter_values[utils.StrToInt(string(values[0]))], string(values[2]))
			}
		}
	}
	for _, filter_id := range filter_ids {
		this.specs = append(this.specs, &ShopProductSpec{wrap: this.wrap, object: &utils.MySql_shop_product_spec{
			A_product_id:   this.object.A_id,
			A_filter_id:    filter_id,
			A_filter_name:  filter_names[filter_id],
			A_filter_value: strings.Join(filter_values[filter_id], ", "),
		}})
	}

	return this
}

func (this *ShopProduct) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *ShopProduct) User() *User {
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

func (this *ShopProduct) Currency() *Currency {
	if this == nil {
		return nil
	}
	if this.currency != nil {
		return this.currency
	}
	this.currency = &Currency{wrap: this.wrap}
	this.currency.loadById(this.object.A_currency)
	return this.currency
}

func (this *ShopProduct) Price() float64 {
	if this == nil {
		return 0
	}
	return this.object.A_price
}

func (this *ShopProduct) PriceFormat(format string) string {
	if this == nil {
		return ""
	}
	return utils.Float64ToStrF(this.object.A_price, format)
}

func (this *ShopProduct) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *ShopProduct) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *ShopProduct) Vendor() string {
	if this == nil {
		return ""
	}
	return this.object.A_vendor
}

func (this *ShopProduct) Quantity() int {
	if this == nil {
		return 0
	}
	return this.object.A_quantity
}

func (this *ShopProduct) Category() *ShopCategory {
	if this == nil {
		return nil
	}
	if this.category != nil {
		return this.category
	}
	this.category = &ShopCategory{wrap: this.wrap}
	this.category.loadById(this.object.A_category)
	return this.category
}

func (this *ShopProduct) Briefly() template.HTML {
	if this == nil {
		return template.HTML("")
	}
	return template.HTML(this.object.A_briefly)
}

func (this *ShopProduct) Content() template.HTML {
	if this == nil {
		return template.HTML("")
	}
	return template.HTML(this.object.A_content)
}

func (this *ShopProduct) DateTimeUnix() int {
	if this == nil {
		return 0
	}
	return this.object.A_datetime
}

func (this *ShopProduct) DateTimeFormat(format string) string {
	if this == nil {
		return ""
	}
	return time.Unix(int64(this.object.A_datetime), 0).Format(format)
}

func (this *ShopProduct) Active() bool {
	if this == nil {
		return false
	}
	return this.object.A_active > 0
}

func (this *ShopProduct) Permalink() string {
	if this == nil {
		return ""
	}
	return "/shop/" + this.object.A_alias + "/"
}

func (this *ShopProduct) Image() *ShopProductImage {
	if this == nil {
		return nil
	}
	if len(this.images) <= 0 {
		return nil
	}
	return this.images[0]
}

func (this *ShopProduct) HaveImages() bool {
	if this == nil {
		return false
	}
	if len(this.images) <= 0 {
		return false
	}
	return true
}

func (this *ShopProduct) Images() []*ShopProductImage {
	if this == nil {
		return []*ShopProductImage{}
	}
	return this.images
}

func (this *ShopProduct) ImagesCount() int {
	if this == nil {
		return 0
	}
	return len(this.images)
}

func (this *ShopProduct) HaveSpecs() bool {
	if this == nil {
		return false
	}
	if len(this.specs) <= 0 {
		return false
	}
	return true
}

func (this *ShopProduct) Specs() []*ShopProductSpec {
	if this == nil {
		return []*ShopProductSpec{}
	}
	return this.specs
}

func (this *ShopProduct) SpecsCount() int {
	if this == nil {
		return 0
	}
	return len(this.specs)
}
