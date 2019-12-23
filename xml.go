package main

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper/config"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-worker/worker"
)

func xml_generator(www_dir string, mp *mysqlpool.MySqlPool) *worker.Worker {
	return worker.New(func(ctx context.Context, w *worker.Worker, o *[]worker.Iface) {
		if www_dir, ok := (*o)[0].(string); ok {
			if mp, ok := (*o)[1].(*mysqlpool.MySqlPool); ok {
				xml_loop(ctx, www_dir, mp)
			}
		}
		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
			return
		}
	}, &[]worker.Iface{
		www_dir,
		mp,
	})
}

func xml_loop(ctx context.Context, www_dir string, mp *mysqlpool.MySqlPool) {
	dirs, err := ioutil.ReadDir(www_dir)
	if err == nil {
		for _, dir := range dirs {
			select {
			case <-ctx.Done():
				return
			default:
				if mp != nil {
					target_dir := strings.Join([]string{www_dir, dir.Name()}, string(os.PathSeparator))
					if utils.IsDirExists(target_dir) {
						xml_detect(ctx, target_dir, dir.Name(), mp)
					}
				}
			}
		}
	}
}

func xml_detect(ctx context.Context, dir, host string, mp *mysqlpool.MySqlPool) {
	db := mp.Get(host)
	if db != nil {
		trigger := strings.Join([]string{dir, "tmp", "trigger.xml.run"}, string(os.PathSeparator))
		if utils.IsFileExists(trigger) {
			if err := db.Ping(); err == nil {
				xml_create(ctx, dir, host, db)
				os.Remove(trigger)
			}
		}
	}
}

func xml_create(ctx context.Context, dir, host string, db *sqlw.DB) {
	conf := config.ConfigNew()
	if err := conf.ConfigRead(strings.Join([]string{dir, "config", "config.json"}, string(os.PathSeparator))); err == nil {
		if (*conf).API.XML.Enabled == 1 {
			if file, err := os.Create(strings.Join([]string{dir, "htdocs", "products.xml"}, string(os.PathSeparator))); err == nil {
				file.Write([]byte(xml_generate(ctx, db, conf)))
				file.Close()
			} else {
				fmt.Printf("Xml generation error (file): %v\n", err)
			}
		}
	} else {
		fmt.Printf("Xml generation error (config): %v\n", err)
	}
}

func xml_generate(ctx context.Context, db *sqlw.DB, conf *config.Config) string {
	return `<?xml version="1.0" encoding="UTF-8"?>` +
		`<!DOCTYPE yml_catalog SYSTEM "shops.dtd">` +
		`<yml_catalog date="` + time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04") + `">` +
		`<shop>` +
		`<name>` + html.EscapeString((*conf).API.XML.Name) + `</name>` +
		`<company>` + html.EscapeString((*conf).API.XML.Company) + `</company>` +
		`<url>` + html.EscapeString((*conf).API.XML.Url) + `</url>` +
		`<currencies>` + xml_gen_currencies(ctx, db, conf) + `</currencies>` +
		`<categories>` + xml_gen_categories(ctx, db, conf) + `</categories>` +
		`<offers>` + xml_gen_offers(ctx, db, conf) + `</offers>` +
		`</shop>` +
		`</yml_catalog>`
}

func xml_gen_currencies(ctx context.Context, db *sqlw.DB, conf *config.Config) string {
	result := ``
	rows, err := db.Query(
		`SELECT
			code,
			coefficient
		FROM
			shop_currencies
		ORDER BY
			id ASC
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 2)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<currency id="` + html.EscapeString(string(values[0])) + `" rate="` + html.EscapeString(string(values[1])) + `"/>`
			}
		}
	}
	return result
}

func xml_gen_categories(ctx context.Context, db *sqlw.DB, conf *config.Config) string {
	result := ``
	rows, err := db.Query(
		`SELECT
			data.id,
			data.user,
			data.name,
			data.alias,
			data.lft,
			data.rgt,
			MAX(data.parent_id) AS parent_id
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
			) AS data
		WHERE
			data.id <> data.parent_id
		GROUP BY
			data.id
		ORDER BY
			data.lft ASC
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 7)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				if utils.StrToInt(string(values[6])) > 1 {
					result += `<category id="` + html.EscapeString(string(values[0])) + `" parentId="` + html.EscapeString(string(values[6])) + `">` + html.EscapeString(string(values[2])) + `</category>`
				} else {
					result += `<category id="` + html.EscapeString(string(values[0])) + `">` + html.EscapeString(string(values[2])) + `</category>`
				}
			}
		}
	}
	return result
}

func xml_gen_offers(ctx context.Context, db *sqlw.DB, conf *config.Config) string {
	result := ``
	rows, err := db.Query(
		`SELECT
			shop_products.id,
			shop_currencies.code,
			shop_products.price,
			shop_products.name,
			shop_products.alias,
			shop_products.vendor,
			shop_products.quantity,
			shop_products.category,
			shop_products.content,
			IFNULL(shop_products.parent_id, 0),
			shop_products.price_old,
			shop_products.price_promo
		FROM
			shop_products
			LEFT JOIN shop_currencies ON shop_currencies.id = shop_products.currency
		WHERE
			shop_products.active = 1 AND
			shop_products.category > 1
		ORDER BY
			shop_products.id
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 12)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<offer id="` + html.EscapeString(string(values[0])) + `" available="true">`
				result += `<url>` + html.EscapeString((*conf).API.XML.Url) + `shop/` + html.EscapeString(string(values[4])) + `/</url>`
				result += `<price>` + utils.Float64ToStrF(utils.StrToFloat64(string(values[2])), "%.2f") + `</price>`
				if utils.StrToFloat64(string(values[10])) > 0 {
					result += `<price_old>` + utils.Float64ToStrF(utils.StrToFloat64(string(values[10])), "%.2f") + `</price_old>`
				}
				if utils.StrToFloat64(string(values[11])) > 0 {
					result += `<price_promo>` + utils.Float64ToStrF(utils.StrToFloat64(string(values[11])), "%.2f") + `</price_promo>`
				}
				result += `<currencyId>` + html.EscapeString(string(values[1])) + `</currencyId>`
				result += `<categoryId>` + html.EscapeString(string(values[7])) + `</categoryId>`
				result += xml_gen_offer_pictures(ctx, db, conf, utils.StrToInt(string(values[0])), utils.StrToInt(string(values[9])))
				result += `<vendor>` + html.EscapeString(string(values[5])) + `</vendor>`
				result += `<stock_quantity>` + html.EscapeString(string(values[6])) + `</stock_quantity>`
				result += `<name>` + html.EscapeString(string(values[3])) + ` ` + html.EscapeString(string(values[0])) + `</name>`
				result += `<description><![CDATA[` + string(values[8]) + `]]></description>`
				result += xml_gen_offer_attributes(ctx, db, conf, utils.StrToInt(string(values[0])))
				result += `</offer>`
			}
		}
	}
	return result
}

func xml_gen_offer_pictures(ctx context.Context, db *sqlw.DB, conf *config.Config, product_id, parent_id int) string {
	result := ``
	if rows, err := db.Query(
		`SELECT
			shop_product_images.product_id,
			shop_product_images.filename
		FROM
			shop_product_images
		WHERE
			shop_product_images.product_id = ?
		ORDER BY
			shop_product_images.ord ASC
		;`,
		product_id,
	); err == nil {
		defer rows.Close()
		values := make([]string, 2)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<picture>` + html.EscapeString((*conf).API.XML.Url) + `products/images/` + html.EscapeString(string(values[0])) + `/` + html.EscapeString(string(values[1])) + `</picture>`
			}
		}
	}

	// Get images from parent
	if result == "" && parent_id > 0 {
		if rows, err := db.Query(
			`SELECT
				shop_product_images.product_id,
				shop_product_images.filename
			FROM
				shop_product_images
			WHERE
				shop_product_images.product_id = ?
			ORDER BY
				shop_product_images.ord ASC
			;`,
			parent_id,
		); err == nil {
			defer rows.Close()
			values := make([]string, 2)
			scan := make([]interface{}, len(values))
			for i := range values {
				scan[i] = &values[i]
			}
			for rows.Next() {
				err = rows.Scan(scan...)
				if err == nil {
					result += `<picture>` + html.EscapeString((*conf).API.XML.Url) + `products/images/` + html.EscapeString(string(values[0])) + `/` + html.EscapeString(string(values[1])) + `</picture>`
				}
			}
		}
	}

	return result
}

func xml_gen_offer_attributes(ctx context.Context, db *sqlw.DB, conf *config.Config, product_id int) string {
	result := ``
	filter_ids := []int{}
	filter_names := map[int]string{}
	filter_values := map[int][]string{}
	rows, err := db.Query(
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
		product_id,
	)
	if err == nil {
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
				filter_names[utils.StrToInt(string(values[0]))] = html.EscapeString(string(values[1]))
				filter_values[utils.StrToInt(string(values[0]))] = append(filter_values[utils.StrToInt(string(values[0]))], string(values[2]))
			}
		}
	}
	for _, filter_id := range filter_ids {
		result += `<param name="` + html.EscapeString(filter_names[filter_id]) + `">` + html.EscapeString(strings.Join(filter_values[filter_id], ", ")) + `</param>`
	}
	return result
}
