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
			{Mount: "default", Name: "Robots.txt", Show: true, Icon: assets.SysSvgIconBug},
			{Mount: "pagination", Name: "Pagination", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "thumbnails", Name: "Thumbnails", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "api", Name: "API", Show: true, Icon: assets.SysSvgIconList},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""

		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
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
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price">Shop thumbnail 1</label>` +
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
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price">Shop thumbnail 2</label>` +
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
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price">Shop thumbnail 3</label>` +
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
						resize_list += `</select>`
						return `<div class="form-group n3">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_price">Shop thumbnail full</label>` +
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
