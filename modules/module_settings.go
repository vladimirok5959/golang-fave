package modules

import (
	"html"
	"io/ioutil"
	"os"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_Settings() *Module {
	return this.newModule(MInfo{
		WantDB: false,
		Mount:  "settings",
		Name:   "Settings",
		Order:  801,
		System: true,
		Icon:   assets.SysSvgIconGear,
		Sub: &[]MISub{
			{Mount: "default", Name: "General", Show: true, Icon: assets.SysSvgIconGear},
			{Mount: "robots-txt", Name: "Robots.txt", Show: true, Icon: assets.SysSvgIconBug},
			{Mount: "pagination", Name: "Pagination", Show: true, Icon: assets.SysSvgIconPagination},
			{Mount: "thumbnails", Name: "Thumbnails", Show: true, Icon: assets.SysSvgIconThumbnails},
			{Mount: "smtp", Name: "SMTP", Show: true, Icon: assets.SysSvgIconEmail},
			{Mount: "shop", Name: "Shop", Show: true, Icon: assets.SysSvgIconShop},
			{Mount: "api", Name: "API", Show: true, Icon: assets.SysSvgIconApi},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""

		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "General"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-general",
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						modules_list := ``
						modules_list += `<select class="form-control" id="lbl_module-at-home" name="module-at-home">`
						modules_list += `<option value="0"`
						if (*wrap.Config).Engine.MainModule == 0 {
							modules_list += ` selected`
						}
						modules_list += `>Pages</option>`
						modules_list += `<option value="1"`
						if (*wrap.Config).Engine.MainModule == 1 {
							modules_list += ` selected`
						}
						modules_list += `>Blog</option>`
						modules_list += `<option value="2"`
						if (*wrap.Config).Engine.MainModule == 2 {
							modules_list += ` selected`
						}
						modules_list += `>Shop</option>`
						modules_list += `</select>`

						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_module-at-home">Module at home page</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							modules_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "robots-txt" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Robots.txt"},
			})

			fcont := []byte(``)
			fcont, _ = ioutil.ReadFile(wrap.DTemplate + string(os.PathSeparator) + "robots.txt")

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-robots-txt",
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group last"><div class="row"><div class="col-12"><textarea class="form-control autosize" id="lbl_content" name="content" placeholder="" autocomplete="off">` + html.EscapeString(string(fcont)) + `</textarea></div></div></div>`
					},
				},
				{
					Kind: builder.DFKSubmit,
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="row d-lg-none"><div class="col-12"><div class="pt-3"><button type="submit" class="btn btn-primary" data-target="add-edit-button">Save</button></div></div></div>`
					},
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "pagination" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Pagination"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-pagination",
				},
				{
					Kind:     builder.DFKNumber,
					Caption:  "Blog main page",
					Name:     "blog-index",
					Min:      "1",
					Max:      "100",
					Required: true,
					Value:    utils.IntToStr((*wrap.Config).Blog.Pagination.Index),
				},
				{
					Kind:     builder.DFKNumber,
					Caption:  "Blog category page",
					Name:     "blog-category",
					Min:      "1",
					Max:      "100",
					Required: true,
					Value:    utils.IntToStr((*wrap.Config).Blog.Pagination.Category),
				},
				{
					Kind:    builder.DFKText,
					Caption: "",
					Name:    "",
					Value:   "",
					CallBack: func(field *builder.DataFormField) string {
						return `<hr>`
					},
				},
				{
					Kind:     builder.DFKNumber,
					Caption:  "Shop main page",
					Name:     "shop-index",
					Min:      "1",
					Max:      "100",
					Required: true,
					Value:    utils.IntToStr((*wrap.Config).Shop.Pagination.Index),
				},
				{
					Kind:     builder.DFKNumber,
					Caption:  "Shop category page",
					Name:     "shop-category",
					Min:      "1",
					Max:      "100",
					Required: true,
					Value:    utils.IntToStr((*wrap.Config).Shop.Pagination.Category),
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "thumbnails" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Thumbnails"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-thumbnails",
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						resize_list := ``
						resize_list += `<select class="form-control" name="shop-thumbnail-r-1">`
						resize_list += `<option value="0"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail1[2] == 0 {
							resize_list += ` selected`
						}
						resize_list += `>Crop</option>`
						resize_list += `<option value="1"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail1[2] == 1 {
							resize_list += ` selected`
						}
						resize_list += `>Resize</option>`
						resize_list += `<option value="2"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail1[2] == 2 {
							resize_list += ` selected`
						}
						resize_list += `>Fit into size</option>`
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Shop thumbnail 1</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-w-1" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail1[0]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-h-1" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail1[1]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-6">` +
							resize_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						resize_list := ``
						resize_list += `<select class="form-control" name="shop-thumbnail-r-2">`
						resize_list += `<option value="0"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail2[2] == 0 {
							resize_list += ` selected`
						}
						resize_list += `>Crop</option>`
						resize_list += `<option value="1"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail2[2] == 1 {
							resize_list += ` selected`
						}
						resize_list += `>Resize</option>`
						resize_list += `<option value="2"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail2[2] == 2 {
							resize_list += ` selected`
						}
						resize_list += `>Fit into size</option>`
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Shop thumbnail 2</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-w-2" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail2[0]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-h-2" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail2[1]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-6">` +
							resize_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						resize_list := ``
						resize_list += `<select class="form-control" name="shop-thumbnail-r-3">`
						resize_list += `<option value="0"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail3[2] == 0 {
							resize_list += ` selected`
						}
						resize_list += `>Crop</option>`
						resize_list += `<option value="1"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail3[2] == 1 {
							resize_list += ` selected`
						}
						resize_list += `>Resize</option>`
						resize_list += `<option value="2"`
						if (*wrap.Config).Shop.Thumbnails.Thumbnail3[2] == 2 {
							resize_list += ` selected`
						}
						resize_list += `>Fit into size</option>`
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Shop thumbnail 3</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-w-3" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail3[0]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-h-3" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.Thumbnail3[1]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-6">` +
							resize_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						resize_list := ``
						resize_list += `<select class="form-control" name="shop-thumbnail-r-full">`
						resize_list += `<option value="0"`
						if (*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] == 0 {
							resize_list += ` selected`
						}
						resize_list += `>Crop</option>`
						resize_list += `<option value="1"`
						if (*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] == 1 {
							resize_list += ` selected`
						}
						resize_list += `>Resize</option>`
						resize_list += `<option value="2"`
						if (*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] == 2 {
							resize_list += ` selected`
						}
						resize_list += `>Fit into size</option>`
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label>Shop thumbnail full</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-w-full" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.ThumbnailFull[0]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-3">` +
							`<div><input class="form-control" type="number" name="shop-thumbnail-h-full" value="` + utils.IntToStr((*wrap.Config).Shop.Thumbnails.ThumbnailFull[1]) + `" min="100" max="1000" placeholder="" autocomplete="off" required></div>` +
							`<div class="d-md-none mb-3"></div>` +
							`</div>` +
							`<div class="col-md-6">` +
							resize_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "smtp" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "SMTP"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-smtp",
				},
				{
					Kind:    builder.DFKText,
					Caption: "SMTP server host",
					Name:    "smtp-host",
					Value:   (*wrap.Config).SMTP.Host,
					Hint:    "Example: smtp.gmail.com",
				},
				{
					Kind:     builder.DFKNumber,
					Caption:  "SMTP server port",
					Name:     "smtp-port",
					Min:      "0",
					Max:      "9999",
					Required: true,
					Value:    utils.IntToStr((*wrap.Config).SMTP.Port),
					Hint:     "Example: 587",
				},
				{
					Kind:    builder.DFKText,
					Caption: "SMTP user login",
					Name:    "smtp-login",
					Value:   (*wrap.Config).SMTP.Login,
					Hint:    "Example: example@gmail.com",
				},
				{
					Kind:    builder.DFKPassword,
					Caption: "SMTP user password",
					Name:    "smtp-password",
					Value:   "",
					Hint:    "Leave this field empty if you don't want change password",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Email address for testing",
					Name:    "smtp-test-email",
					Value:   "",
					Hint:    "To this email address will be send test email message if settings are correct",
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "shop" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Shop"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-shop",
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						price_format_list := ``
						price_format_list += `<select class="form-control" id="lbl_price-fomat" name="price-fomat">`
						price_format_list += `<option value="0"`
						if (*wrap.Config).Shop.Price.Format == 0 {
							price_format_list += ` selected`
						}
						price_format_list += `>100</option>`
						price_format_list += `<option value="1"`
						if (*wrap.Config).Shop.Price.Format == 1 {
							price_format_list += ` selected`
						}
						price_format_list += `>100.0</option>`
						price_format_list += `<option value="2"`
						if (*wrap.Config).Shop.Price.Format == 2 {
							price_format_list += ` selected`
						}
						price_format_list += `>100.00</option>`
						price_format_list += `<option value="3"`
						if (*wrap.Config).Shop.Price.Format == 3 {
							price_format_list += ` selected`
						}
						price_format_list += `>100.000</option>`
						price_format_list += `<option value="4"`
						if (*wrap.Config).Shop.Price.Format == 4 {
							price_format_list += ` selected`
						}
						price_format_list += `>100.0000</option>`
						price_format_list += `</select>`

						price_round_list := ``
						price_round_list += `<select class="form-control" id="lbl_price-round" name="price-round">`
						price_round_list += `<option value="0"`
						if (*wrap.Config).Shop.Price.Round == 0 {
							price_round_list += ` selected`
						}
						price_round_list += `>Don't round</option>`
						price_round_list += `<option value="1"`
						if (*wrap.Config).Shop.Price.Round == 1 {
							price_round_list += ` selected`
						}
						price_round_list += `>Round to ceil</option>`
						price_round_list += `<option value="2"`
						if (*wrap.Config).Shop.Price.Round == 2 {
							price_round_list += ` selected`
						}
						price_round_list += `>Round to floor</option>`
						price_round_list += `</select>`

						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price-fomat">Price format</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							price_format_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price-round">Price round</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							price_round_list +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						return `<div><hr></div>` +
							`<div><h4>Order process require fields</h4></div>` +
							`<div><hr></div>`
					},
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Last Name",
					Name:    "require-last-name",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.LastName),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "First Name",
					Name:    "require-first-name",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.FirstName),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Second Name",
					Name:    "require-second-name",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.SecondName),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Mobile Phone",
					Name:    "require-mobile-phone",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.MobilePhone),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Email Address",
					Name:    "require-email-address",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.EmailAddress),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Delivery",
					Name:    "require-delivery",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.Delivery),
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Comment",
					Name:    "require-comment",
					Value:   utils.IntToStr((*wrap.Config).Shop.Orders.RequiredFields.Comment),
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		} else if wrap.CurrSubModule == "api" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "API"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "settings-api",
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "XML enabled",
					Name:    "xml-enabled",
					Value:   utils.IntToStr((*wrap.Config).API.XML.Enabled),
					Hint:    "XML: <a href=\"/api/products/\" target=\"_blank\">/api/products/</a>",
				},
				{
					Kind:    builder.DFKText,
					Caption: "XML name",
					Name:    "xml-name",
					Value:   (*wrap.Config).API.XML.Name,
				},
				{
					Kind:    builder.DFKText,
					Caption: "XML company",
					Name:    "xml-company",
					Value:   (*wrap.Config).API.XML.Company,
				},
				{
					Kind:    builder.DFKText,
					Caption: "XML url",
					Name:    "xml-url",
					Value:   (*wrap.Config).API.XML.Url,
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  "Save",
					Target: "add-edit-button",
				},
			})

			sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
