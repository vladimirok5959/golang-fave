package modules

import (
	"net/http"
	"os"
	"strings"

	"golang-fave/assets"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

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
		if len(wrap.UrlArgs) == 2 && wrap.UrlArgs[0] == "api" && wrap.UrlArgs[1] == "products" {
			if (*wrap.Config).API.XML.Enabled == 1 {
				// Fix url
				if wrap.R.URL.Path[len(wrap.R.URL.Path)-1] != '/' {
					http.Redirect(wrap.W, wrap.R, wrap.R.URL.Path+"/"+utils.ExtractGetParams(wrap.R.RequestURI), 301)
					return
				}

				target_file := wrap.DHtdocs + string(os.PathSeparator) + "products.xml"
				if !utils.IsFileExists(target_file) {
					data := []byte(this.api_GenerateEmptyXml(wrap))

					// Make empty file
					if file, err := os.Create(target_file); err == nil {
						file.Write(data)
					}

					// Make regular XML
					data = []byte(this.api_GenerateXml(wrap))

					// Save file
					wrap.RemoveProductXmlCacheFile()
					if file, err := os.Create(target_file); err == nil {
						file.Write(data)
					}

					wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					wrap.W.Header().Set("Content-Type", "text/xml; charset=utf-8")
					wrap.W.WriteHeader(http.StatusOK)
					wrap.W.Write(data)
				} else {
					http.ServeFile(wrap.W, wrap.R, target_file)
				}
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
			wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
			return
		}
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		// No any page for back-end
		return "", "", ""
	})
}

func (this *Modules) RegisterModule_ApiProducts() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "products",
		Name:   "Api Products",
		Order:  804,
		System: true,
		Icon:   assets.SysSvgIconPage,
		Sub:    &[]MISub{},
	}, func(wrap *wrapper.Wrapper) {
		if len(wrap.UrlArgs) == 4 && wrap.UrlArgs[0] == "products" && wrap.UrlArgs[1] == "images" && utils.IsNumeric(wrap.UrlArgs[2]) && wrap.UrlArgs[3] != "" {
			thumb_type := ""
			file_name := ""

			if strings.HasPrefix(wrap.UrlArgs[3], "thumb-0-") {
				thumb_type = "thumb-0"
				file_name = wrap.UrlArgs[3][len(thumb_type)+1:]
			} else if strings.HasPrefix(wrap.UrlArgs[3], "thumb-1-") {
				thumb_type = "thumb-1"
				file_name = wrap.UrlArgs[3][len(thumb_type)+1:]
			} else if strings.HasPrefix(wrap.UrlArgs[3], "thumb-2-") {
				thumb_type = "thumb-2"
				file_name = wrap.UrlArgs[3][len(thumb_type)+1:]
			} else if strings.HasPrefix(wrap.UrlArgs[3], "thumb-3-") {
				thumb_type = "thumb-3"
				file_name = wrap.UrlArgs[3][len(thumb_type)+1:]
			} else if strings.HasPrefix(wrap.UrlArgs[3], "thumb-full-") {
				thumb_type = "thumb-full"
				file_name = wrap.UrlArgs[3][len(thumb_type)+1:]
			}

			if !(thumb_type == "" && file_name == "") {
				original_file := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + wrap.UrlArgs[2] + string(os.PathSeparator) + file_name
				if !utils.IsFileExists(original_file) {
					// User error 404 page
					wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
					return
				}

				width := (*wrap.Config).Shop.Thumbnails.Thumbnail0[0]
				height := (*wrap.Config).Shop.Thumbnails.Thumbnail0[1]
				resize := false

				if thumb_type == "thumb-1" {
					width = (*wrap.Config).Shop.Thumbnails.Thumbnail1[0]
					height = (*wrap.Config).Shop.Thumbnails.Thumbnail1[1]
					if (*wrap.Config).Shop.Thumbnails.Thumbnail1[2] == 1 {
						resize = true
					}
				} else if thumb_type == "thumb-2" {
					width = (*wrap.Config).Shop.Thumbnails.Thumbnail2[0]
					height = (*wrap.Config).Shop.Thumbnails.Thumbnail2[1]
					if (*wrap.Config).Shop.Thumbnails.Thumbnail2[2] == 1 {
						resize = true
					}
				} else if thumb_type == "thumb-3" {
					width = (*wrap.Config).Shop.Thumbnails.Thumbnail3[0]
					height = (*wrap.Config).Shop.Thumbnails.Thumbnail3[1]
					if (*wrap.Config).Shop.Thumbnails.Thumbnail3[2] == 1 {
						resize = true
					}
				} else if thumb_type == "thumb-full" {
					width = (*wrap.Config).Shop.Thumbnails.ThumbnailFull[0]
					height = (*wrap.Config).Shop.Thumbnails.ThumbnailFull[1]
					if (*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] == 1 {
						resize = true
					}
				}

				target_file := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + wrap.UrlArgs[2] + string(os.PathSeparator) + thumb_type + "-" + file_name
				if !utils.IsFileExists(target_file) {
					data, ok, ext, err := this.api_GenerateImage(wrap, width, height, resize, original_file)
					if err != nil {
						// System error 500
						utils.SystemErrorPageEngine(wrap.W, err)
						return
					}

					if !ok {
						// User error 404 page
						wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
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
			} else {
				// User error 404 page
				wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
				return
			}
		} else {
			// User error 404 page
			wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
			return
		}
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		// No any page for back-end
		return "", "", ""
	})
}
