package modules

import (
	"bufio"
	"bytes"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang-fave/assets"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"

	"github.com/disintegration/imaging"
)

func (this *Modules) api_GenerateImage(wrap *wrapper.Wrapper, width, height, color int, filename string) ([]byte, bool, string, error) {
	file_ext := ""
	if strings.ToLower(filepath.Ext(filename)) == ".png" {
		file_ext = "image/png"
	} else if strings.ToLower(filepath.Ext(filename)) == ".jpg" {
		file_ext = "image/jpeg"
	} else if strings.ToLower(filepath.Ext(filename)) == ".jpeg" {
		file_ext = "image/jpeg"
	}

	src, err := imaging.Open(filename)
	if err != nil {
		return []byte(""), false, file_ext, err
	}

	src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
	// src = imaging.Fit(src, width, height, imaging.Lanczos)

	var out_bytes bytes.Buffer
	out := bufio.NewWriter(&out_bytes)

	if file_ext == "image/png" {
		imaging.Encode(out, src, imaging.PNG)
	} else if file_ext == "image/jpeg" {
		imaging.Encode(out, src, imaging.JPEG)
	} else {
		return []byte(""), false, file_ext, nil
	}

	return out_bytes.Bytes(), true, file_ext, nil
}

func (this *Modules) api_GenerateXmlCurrencies(wrap *wrapper.Wrapper) string {
	result := ``
	rows, err := wrap.DB.Query(
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

func (this *Modules) api_GenerateXmlCategories(wrap *wrapper.Wrapper) string {
	result := ``
	rows, err := wrap.DB.Query(
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

func (this *Modules) api_GenerateXmlOfferPictures(wrap *wrapper.Wrapper, product_id int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			shop_product_images.product_id,
			shop_product_images.filename
		FROM
			shop_product_images
		WHERE
			shop_product_images.product_id = ?
		;`,
		product_id,
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
				result += `<picture>` + html.EscapeString((*wrap.Config).API.XML.Url) + `products/images/` + html.EscapeString(string(values[0])) + `/` + html.EscapeString(string(values[1])) + `</picture>`
			}
		}
	}
	return result
}

func (this *Modules) api_GenerateXmlOfferAttributes(wrap *wrapper.Wrapper, product_id int) string {
	result := ``
	filter_ids := []int{}
	filter_names := map[int]string{}
	filter_values := map[int][]string{}
	rows, err := wrap.DB.Query(
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

func (this *Modules) api_GenerateXmlOffers(wrap *wrapper.Wrapper) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			shop_products.id,
			shop_currencies.code,
			shop_products.price,
			shop_products.name,
			shop_products.alias,
			shop_products.vendor,
			shop_products.quantity,
			shop_products.category,
			shop_products.content
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
		values := make([]string, 9)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				result += `<offer id="` + html.EscapeString(string(values[0])) + `" available="true">`
				result += `<url>` + html.EscapeString((*wrap.Config).API.XML.Url) + `shop/` + html.EscapeString(string(values[4])) + `/</url>`
				result += `<price>` + utils.Float64ToStrF(utils.StrToFloat64(string(values[2])), "%.2f") + `</price>`
				result += `<currencyId>` + html.EscapeString(string(values[1])) + `</currencyId>`
				result += `<categoryId>` + html.EscapeString(string(values[7])) + `</categoryId>`
				result += this.api_GenerateXmlOfferPictures(wrap, utils.StrToInt(string(values[0])))
				result += `<vendor>` + html.EscapeString(string(values[5])) + `</vendor>`
				result += `<stock_quantity>` + html.EscapeString(string(values[6])) + `</stock_quantity>`
				result += `<name>` + html.EscapeString(string(values[3])) + `</name>`
				result += `<description><![CDATA[` + string(values[8]) + `]]></description>`
				result += this.api_GenerateXmlOfferAttributes(wrap, utils.StrToInt(string(values[0])))
				result += `</offer>`
			}
		}
	}
	return result
}

func (this *Modules) api_GenerateXml(wrap *wrapper.Wrapper) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE yml_catalog SYSTEM "shops.dtd">
<yml_catalog date="` + time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04") + `">
	<shop>
		<name>` + html.EscapeString((*wrap.Config).API.XML.Name) + `</name>
		<company>` + html.EscapeString((*wrap.Config).API.XML.Company) + `</company>
		<url>` + html.EscapeString((*wrap.Config).API.XML.Url) + `</url>
		<currencies>` + this.api_GenerateXmlCurrencies(wrap) + `</currencies>
		<categories>` + this.api_GenerateXmlCategories(wrap) + `</categories>
		<offers>` + this.api_GenerateXmlOffers(wrap) + `</offers>
	</shop>
</yml_catalog>`
}

func (this *Modules) RegisterModule_Api() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "api",
		Name:   "Api",
		Order:  803,
		System: true,
		Icon:   assets.SysSvgIconPage,
		Sub:    &[]MISub{},
	}, func(wrap *wrapper.Wrapper) {
		if len(wrap.UrlArgs) == 5 && wrap.UrlArgs[0] == "api" && wrap.UrlArgs[1] == "product-image" && (wrap.UrlArgs[2] == "thumb-0" || wrap.UrlArgs[2] == "thumb-1" || wrap.UrlArgs[2] == "thumb-2" || wrap.UrlArgs[2] == "thumb-3") {
			thumb_type := wrap.UrlArgs[2]
			product_id := wrap.UrlArgs[3]
			file_name := wrap.UrlArgs[4]

			original_file := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + product_id + string(os.PathSeparator) + file_name
			if !utils.IsFileExists(original_file) {
				// User error 404 page
				wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true), http.StatusNotFound)
				return
			}

			width := (*wrap.Config).Shop.Thumbnails.Thumbnail0[0]
			height := (*wrap.Config).Shop.Thumbnails.Thumbnail0[1]

			if thumb_type == "thumb-1" {
				width = (*wrap.Config).Shop.Thumbnails.Thumbnail1[0]
				height = (*wrap.Config).Shop.Thumbnails.Thumbnail1[1]
			} else if thumb_type == "thumb-2" {
				width = (*wrap.Config).Shop.Thumbnails.Thumbnail2[0]
				height = (*wrap.Config).Shop.Thumbnails.Thumbnail2[1]
			} else if thumb_type == "thumb-3" {
				width = (*wrap.Config).Shop.Thumbnails.Thumbnail3[0]
				height = (*wrap.Config).Shop.Thumbnails.Thumbnail3[1]
			}

			target_file := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + product_id + string(os.PathSeparator) + thumb_type + "-" + file_name
			if !utils.IsFileExists(target_file) {
				data, ok, ext, err := this.api_GenerateImage(wrap, width, height, 0, original_file)
				if err != nil {
					// System error 500
					utils.SystemErrorPageEngine(wrap.W, err)
					return
				}

				if !ok {
					// User error 404 page
					wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true), http.StatusNotFound)
					return
				}

				// Save file
				if file, err := os.Create(target_file); err == nil {
					file.Write(data)
				}

				wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				wrap.W.Header().Set("Content-Type", ext)
				wrap.W.Write(data)
			} else {
				http.ServeFile(wrap.W, wrap.R, target_file)
			}
		} else if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "api" && wrap.UrlArgs[1] == "products" {
			if (*wrap.Config).API.XML.Enabled == 1 {
				// Fix url
				if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
					http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
					return
				}

				// XML
				wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				wrap.W.Header().Set("Content-Type", "text/xml; charset=utf-8")
				wrap.W.WriteHeader(http.StatusOK)
				wrap.W.Write([]byte(this.api_GenerateXml(wrap)))
			} else {
				wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				wrap.W.WriteHeader(http.StatusNotFound)
				wrap.W.Write([]byte("Disabled!"))
			}
		} else if len(wrap.UrlArgs) == 1 {
			// Fix url
			if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
				http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
				return
			}

			// Some info
			wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			wrap.W.WriteHeader(http.StatusOK)
			wrap.W.Write([]byte("Fave engine API mount point!"))
		} else {
			// User error 404 page
			wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true), http.StatusNotFound)
			return
		}
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		// No any page for back-end
		return "", "", ""
	})
}
