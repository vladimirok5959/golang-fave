package fetdata

import (
	"bytes"
	"html"
	"html/template"
	"os"
	"time"

	"golang-fave/consts"
	"golang-fave/engine/basket"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	wrap  *wrapper.Wrapper
	is404 bool

	Page *Page
	Blog *Blog
	Shop *Shop
}

func New(wrap *wrapper.Wrapper, is404 bool, drow interface{}, duser *utils.MySql_user) *FERData {
	var d_Page *Page
	var d_Blog *Blog
	var d_Shop *Shop

	var preUser *User
	if duser != nil {
		preUser = &User{wrap: wrap, object: duser}
	}

	if wrap.CurrModule == "index" {
		if o, ok := drow.(*utils.MySql_page); ok {
			d_Page = &Page{wrap: wrap, object: o, user: preUser}
		}
	} else if wrap.CurrModule == "blog" {
		if len(wrap.UrlArgs) == 3 && wrap.UrlArgs[0] == "blog" && wrap.UrlArgs[1] == "category" && wrap.UrlArgs[2] != "" {
			if o, ok := drow.(*utils.MySql_blog_category); ok {
				d_Blog = &Blog{wrap: wrap, category: (&BlogCategory{wrap: wrap, object: o, user: preUser}).load(nil)}
				d_Blog.load()
			}
		} else if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "blog" && wrap.UrlArgs[1] != "" {
			if o, ok := drow.(*utils.MySql_blog_post); ok {
				d_Blog = &Blog{wrap: wrap, post: (&BlogPost{wrap: wrap, object: o, user: preUser}).load()}
			}
		} else {
			d_Blog = &Blog{wrap: wrap}
			d_Blog.load()
		}
	} else if wrap.CurrModule == "shop" {
		if len(wrap.UrlArgs) == 3 && wrap.UrlArgs[0] == "shop" && wrap.UrlArgs[1] == "category" && wrap.UrlArgs[2] != "" {
			if o, ok := drow.(*utils.MySql_shop_category); ok {
				d_Shop = &Shop{wrap: wrap, category: (&ShopCategory{wrap: wrap, object: o, user: preUser}).load(nil)}
				d_Shop.load()
			}
		} else if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "shop" && wrap.UrlArgs[1] != "" {
			if o, ok := drow.(*utils.MySql_shop_product); ok {
				d_Shop = &Shop{wrap: wrap, product: (&ShopProduct{wrap: wrap, object: o, user: preUser}).load()}
			}
		} else {
			d_Shop = &Shop{wrap: wrap}
			d_Shop.load()
		}
	}

	if d_Blog == nil {
		d_Blog = &Blog{wrap: wrap}
	}

	if d_Shop == nil {
		d_Shop = &Shop{wrap: wrap}
	}

	fer := &FERData{
		wrap:  wrap,
		is404: is404,
		Page:  d_Page,
		Blog:  d_Blog,
		Shop:  d_Shop,
	}

	return fer
}

func (this *FERData) RequestURI() string {
	return this.wrap.R.RequestURI
}

func (this *FERData) RequestURL() string {
	return this.wrap.R.URL.Path
}

func (this *FERData) RequestGET() string {
	return utils.ExtractGetParams(this.wrap.R.RequestURI)
}

func (this *FERData) Module() string {
	if this.is404 {
		return "404"
	}
	var mod string
	if this.wrap.CurrModule == "index" {
		mod = "index"
	} else if this.wrap.CurrModule == "blog" {
		if len(this.wrap.UrlArgs) == 3 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] == "category" && this.wrap.UrlArgs[2] != "" {
			mod = "blog-category"
		} else if len(this.wrap.UrlArgs) == 2 && this.wrap.UrlArgs[0] == "blog" && this.wrap.UrlArgs[1] != "" {
			mod = "blog-post"
		} else {
			mod = "blog"
		}
	} else if this.wrap.CurrModule == "shop" {
		if len(this.wrap.UrlArgs) == 3 && this.wrap.UrlArgs[0] == "shop" && this.wrap.UrlArgs[1] == "category" && this.wrap.UrlArgs[2] != "" {
			mod = "shop-category"
		} else if len(this.wrap.UrlArgs) == 2 && this.wrap.UrlArgs[0] == "shop" && this.wrap.UrlArgs[1] != "" {
			mod = "shop-product"
		} else {
			mod = "shop"
		}
	}
	return mod
}

func (this *FERData) DateTimeUnix() int {
	return int(time.Now().Unix())
}

func (this *FERData) DateTimeFormat(format string) string {
	return time.Unix(int64(time.Now().Unix()), 0).Format(format)
}

func (this *FERData) EscapeString(str string) string {
	return html.EscapeString(str)
}

func (this *FERData) cachedBlock(block string) template.HTML {
	tmpl, err := template.New(block + ".html").Funcs(utils.TemplateAdditionalFuncs()).ParseFiles(
		this.wrap.DTemplate + string(os.PathSeparator) + block + ".html",
	)
	if err != nil {
		return template.HTML(err.Error())
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, consts.TmplData{
		System: utils.GetTmplSystemData("", ""),
		Data:   this,
	})
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(string(tpl.Bytes()))
}

func (this *FERData) CachedBlock1() template.HTML {
	if data, ok := this.wrap.GetBlock1(); ok {
		return data
	}
	data := this.cachedBlock("cached-block-1")
	this.wrap.SetBlock1(data)
	return data
}

func (this *FERData) CachedBlock2() template.HTML {
	if data, ok := this.wrap.GetBlock2(); ok {
		return data
	}
	data := this.cachedBlock("cached-block-2")
	this.wrap.SetBlock2(data)
	return data
}

func (this *FERData) CachedBlock3() template.HTML {
	if data, ok := this.wrap.GetBlock3(); ok {
		return data
	}
	data := this.cachedBlock("cached-block-3")
	this.wrap.SetBlock3(data)
	return data
}

func (this *FERData) CachedBlock4() template.HTML {
	if data, ok := this.wrap.GetBlock4(); ok {
		return data
	}
	data := this.cachedBlock("cached-block-4")
	this.wrap.SetBlock4(data)
	return data
}

func (this *FERData) CachedBlock5() template.HTML {
	if data, ok := this.wrap.GetBlock5(); ok {
		return data
	}
	data := this.cachedBlock("cached-block-5")
	this.wrap.SetBlock5(data)
	return data
}

func (this *FERData) ShopBasketProductsCount() int {
	return this.wrap.ShopBasket.ProductsCount(&basket.SBParam{
		R:         this.wrap.R,
		DB:        this.wrap.DB,
		Host:      this.wrap.CurrHost,
		SessionId: this.wrap.GetSessionId(),
	})
}
