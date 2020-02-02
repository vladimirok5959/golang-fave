package workers

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang-fave/engine/config"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"

	"github.com/vladimirok5959/golang-worker/worker"
)

func XmlGenerator(www_dir string, mp *mysqlpool.MySqlPool) *worker.Worker {
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
			if err := db.Ping(ctx); err == nil {
				xml_create(ctx, dir, host, trigger, db)
			}
		}
	}
}

func xml_create(ctx context.Context, dir, host, trigger string, db *sqlw.DB) {
	conf := config.ConfigNew()
	if err := conf.ConfigRead(strings.Join([]string{dir, "config", "config.json"}, string(os.PathSeparator))); err == nil {
		if (*conf).API.XML.Enabled == 1 {
			if file, err := os.Create(strings.Join([]string{dir, "htdocs", "products.xml"}, string(os.PathSeparator))); err == nil {
				if content, err := xml_generate(ctx, db, conf); err == nil {
					if _, err := file.Write([]byte(content)); err == nil {
						os.Remove(trigger)
					}
				}
				file.Close()
			} else {
				fmt.Printf("Xml generation error (file): %v\n", err)
			}
		}
	} else {
		fmt.Printf("Xml generation error (config): %v\n", err)
	}
}

func xml_generate(ctx context.Context, db *sqlw.DB, conf *config.Config) (string, error) {
	content := ""

	var currencies string
	var categories string
	var offers string
	var err error

	if currencies, err = xml_gen_currencies(ctx, db, conf); err != nil {
		return content, err
	}

	if categories, err = xml_gen_categories(ctx, db, conf); err != nil {
		return content, err
	}

	if offers, err = xml_gen_offers(ctx, db, conf); err != nil {
		return content, err
	}

	return `<?xml version="1.0" encoding="UTF-8"?>` +
			`<!DOCTYPE yml_catalog SYSTEM "shops.dtd">` +
			`<yml_catalog date="` + time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04") + `">` +
			`<shop>` +
			`<name>` + html.EscapeString((*conf).API.XML.Name) + `</name>` +
			`<company>` + html.EscapeString((*conf).API.XML.Company) + `</company>` +
			`<url>` + html.EscapeString((*conf).API.XML.Url) + `</url>` +
			`<currencies>` + currencies + `</currencies>` +
			`<categories>` + categories + `</categories>` +
			`<offers>` + offers + `</offers>` +
			`</shop>` +
			`</yml_catalog>`,
		nil
}

func xml_gen_currencies(ctx context.Context, db *sqlw.DB, conf *config.Config) (string, error) {
	result := ``
	rows, err := db.Query(
		ctx,
		`SELECT
			code,
			coefficient
		FROM
			fave_shop_currencies
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
	return result, nil
}

func xml_gen_categories(ctx context.Context, db *sqlw.DB, conf *config.Config) (string, error) {
	result := ``
	rows, err := db.Query(
		ctx,
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
					fave_shop_cats AS node,
					fave_shop_cats AS parent
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
	return result, nil
}

func xml_gen_offers(ctx context.Context, db *sqlw.DB, conf *config.Config) (string, error) {
	result := ``
	rows, err := db.Query(
		ctx,
		`SELECT
			fave_shop_products.id,
			fave_shop_currencies.code,
			fave_shop_products.price,
			fave_shop_products.name,
			fave_shop_products.alias,
			fave_shop_products.vendor,
			fave_shop_products.quantity,
			fave_shop_products.category,
			fave_shop_products.content,
			IFNULL(fave_shop_products.parent_id, 0),
			fave_shop_products.price_old,
			fave_shop_products.price_promo
		FROM
			fave_shop_products
			LEFT JOIN fave_shop_currencies ON fave_shop_currencies.id = fave_shop_products.currency
		WHERE
			fave_shop_products.active = 1 AND
			fave_shop_products.category > 1
		ORDER BY
			fave_shop_products.id
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
	return result, nil
}

func xml_gen_offer_pictures(ctx context.Context, db *sqlw.DB, conf *config.Config, product_id, parent_id int) string {
	result := ``
	if rows, err := db.Query(
		ctx,
		`SELECT
			fave_shop_product_images.product_id,
			fave_shop_product_images.filename
		FROM
			fave_shop_product_images
		WHERE
			fave_shop_product_images.product_id = ?
		ORDER BY
			fave_shop_product_images.ord ASC
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
			ctx,
			`SELECT
				fave_shop_product_images.product_id,
				fave_shop_product_images.filename
			FROM
				fave_shop_product_images
			WHERE
				fave_shop_product_images.product_id = ?
			ORDER BY
				fave_shop_product_images.ord ASC
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
		ctx,
		`SELECT
			fave_shop_filters.id,
			fave_shop_filters.filter,
			fave_shop_filters_values.name
		FROM
			fave_shop_filter_product_values
			LEFT JOIN fave_shop_filters_values ON fave_shop_filters_values.id = fave_shop_filter_product_values.filter_value_id
			LEFT JOIN fave_shop_filters ON fave_shop_filters.id = fave_shop_filters_values.filter_id
		WHERE
			fave_shop_filter_product_values.product_id = ?
		ORDER BY
			fave_shop_filters.filter ASC,
			fave_shop_filters_values.name ASC
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
