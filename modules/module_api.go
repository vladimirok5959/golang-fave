package modules

import (
	"html"
	"net/http"
	"time"

	"golang-fave/assets"
	// "golang-fave/consts"
	// "golang-fave/engine/builder"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

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
	/*
		<category id="2">Женская одежда</category>
		<category id="261" parentId="2">Платья</category>
		<category id="3">Мужская одежда</category>
		<category id="391" parentId="3">Куртки</category>
	*/
	return ``
}

func (this *Modules) api_GenerateXmlOffers(wrap *wrapper.Wrapper) string {
	/*
		<offer id="19305" available="true">
			<url>http://abc.ua/catalog/muzhskaya_odezhda/kurtki/kurtkabx.html</url>
			<price>4499</price>
			<currencyId>UAH</currencyId>
			<categoryId>391</categoryId>
			<picture>http://abc.ua/upload/iblock/a53/a5391cddb40be91705.jpg</picture>
			<picture>http://abc.ua/upload/iblock/9d0/9d06805d219fb525fc.jpg</picture>
			<picture>http://abc.ua/upload/iblock/93d/93de38537e1cc1f8f2.jpg</picture>
			<vendor>Abc clothes</vendor>
			<stock_quantity>100</stock_quantity>
			<name>Куртка Abc clothes Scoperandom-HH XL Черная (1323280942900)</name>
			<description><![CDATA[<p>Одежда<b>Abc clothes</b> способствует развитию функций головного мозга за счет поощрения мелкой моторики.</p><p>В Abc <b>New Collection</b> будет особенно удобно лазать, прыгать, бегать.</p><p>За счет своей универсальноcти и многофункциональности, <b>Abc clothes</b> отлично подходит:</p><ul><li><b>Для весны</b></li><li><b>Для лета</b></li><li><b>Для ранней осени</b> </li></ul><br><p><b>Состав:</b><br>• 92% полиэстер, 8% эластан, нетоксичность подтверждена лабораторно.</p><p><b>Вес:</b> 305 г</p>]]></description>
			<param name="Вид">Куртка</param>
			<param name="Размер">XL</param>
			<param name="Сезон">Весна-Осень</param>
			<param name="Категория">Мужская</param>
			<param name="Цвет">Черный</param>
			<param name="Длина">Средней длины</param>
			<param name="Стиль">Повседневный (casual)</param>
			<param name="Особенности">Модель с капюшоном</param>
			<param name="Состав">92% полиэстер, 8% эластан</param>
			<param name="Артикул">58265468</param>
		</offer>
	*/
	return ``
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
		if (*wrap.Config).API.XML.Enabled == 1 {
			if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "api" && wrap.UrlArgs[1] == "products" {
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
