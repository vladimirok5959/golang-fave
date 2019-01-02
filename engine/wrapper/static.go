package wrapper

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	images "golang-fave/engine/wrapper/resources/images"
	others "golang-fave/engine/wrapper/resources/others"
	scripts "golang-fave/engine/wrapper/resources/scripts"
	styles "golang-fave/engine/wrapper/resources/styles"
	templates "golang-fave/engine/wrapper/resources/templates"
)

type tmplDataErrorMsg struct {
	ErrorMessage string
}

func (this *Wrapper) staticResource() bool {
	if this.R.URL.Path == "/assets/sys/styles.css" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "text/css")
		(*this.W).Write(styles.File_assets_sys_styles_css)
		return true
	} else if this.R.URL.Path == "/assets/cp/styles.css" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "text/css")
		(*this.W).Write(styles.File_assets_cp_styles_css)
		return true
	} else if this.R.URL.Path == "/assets/sys/bootstrap.css" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "text/css")
		(*this.W).Write(styles.File_assets_sys_bootstrap_css)
		return true
	} else if this.R.URL.Path == "/assets/sys/jquery.js" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "application/javascript; charset=utf-8")
		(*this.W).Write(scripts.File_assets_sys_jquery_js)
		return true
	} else if this.R.URL.Path == "/assets/sys/popper.js" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "application/javascript; charset=utf-8")
		(*this.W).Write(scripts.File_assets_sys_popper_js)
		return true
	} else if this.R.URL.Path == "/assets/sys/bootstrap.js" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "application/javascript; charset=utf-8")
		(*this.W).Write(scripts.File_assets_sys_bootstrap_js)
		return true
	} else if this.R.URL.Path == "/assets/cp/scripts.js" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "application/javascript; charset=utf-8")
		(*this.W).Write(scripts.File_assets_cp_scripts_js)
		return true
	} else if this.R.URL.Path == "/assets/sys/logo.svg" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "image/svg+xml")
		(*this.W).Write(others.File_assets_sys_logo_svg)
		return true
	} else if this.R.URL.Path == "/assets/sys/bg.png" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "image/png")
		(*this.W).Write(images.File_assets_sys_bg_png)
		return true
	} else if this.R.URL.Path == "/assets/sys/logo.png" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "image/png")
		(*this.W).Write(images.File_assets_sys_logo_png)
		return true
	} else if this.R.URL.Path == "/assets/sys/fave.ico" {
		(*this.W).Header().Set("Cache-Control", "public, max-age=31536000")
		(*this.W).Header().Set("Content-Type", "image/x-icon")
		(*this.W).Write(others.File_assets_sys_fave_ico)
		return true
	}
	return false
}

func (this *Wrapper) staticFile() bool {
	file := this.R.URL.Path
	if file != "/" {
		f, err := os.Open(this.DirVHostHome + "/htdocs" + file)
		if err == nil {
			defer f.Close()
			st, err := os.Stat(this.DirVHostHome + "/htdocs" + file)
			if err != nil {
				return false
			}
			if st.Mode().IsDir() {
				if file[len(file)-1] == '/' {
					fi, err := os.Open(this.DirVHostHome + "/htdocs" + file + "/index.html")
					if err == nil {
						defer fi.Close()
						sti, err := os.Stat(this.DirVHostHome + "/htdocs" + file + "/index.html")
						if err != nil {
							return false
						}
						if sti.Mode().IsDir() {
							return false
						}
						http.ServeFile(*this.W, this.R, this.DirVHostHome+"/htdocs"+file+"/index.html")
						return true
					}
				}
				return false
			}
			http.ServeFile(*this.W, this.R, this.DirVHostHome+"/htdocs"+file)
			return true
		}
	} else {
		f, err := os.Open(this.DirVHostHome + "/htdocs/index.html")
		if err == nil {
			defer f.Close()
			st, err := os.Stat(this.DirVHostHome + "/htdocs/index.html")
			if err != nil {
				return false
			}
			if st.Mode().IsDir() {
				return false
			}
			http.ServeFile(*this.W, this.R, this.DirVHostHome+"/htdocs/index.html")
			return true
		}
	}
	return false
}

func (this *Wrapper) printPageDefault() {
	// Custom page
	f, err := os.Open(this.DirVHostHome + "/htdocs" + "/index.html")
	if err == nil {
		defer f.Close()
		http.ServeFile(*this.W, this.R, this.DirVHostHome+"/htdocs"+"/index.html")
		return
	}

	// Default page
	tmpl, err := template.New("template").Parse(string(templates.PageDefault))
	if err != nil {
		this.printTmplPageError(err)
		return
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.W).Header().Set("Content-Type", "text/html")
	tmpl.Execute(*this.W, TmplDataAll{
		System: this.TmplGetSystemData(),
	})
}

func (this *Wrapper) printPage404() {
	// Custom 404 error page
	f, err := ioutil.ReadFile(this.DirVHostHome + "/htdocs" + "/404.html")
	if err == nil {
		(*this.W).WriteHeader(http.StatusNotFound)
		(*this.W).Header().Set("Content-Type", "text/html")
		(*this.W).Write(f)
		return
	}

	// Default error page
	tmpl, err := template.New("template").Parse(string(templates.PageError404))
	if err != nil {
		this.printTmplPageError(err)
		return
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.W).WriteHeader(http.StatusNotFound)
	(*this.W).Header().Set("Content-Type", "text/html")
	tmpl.Execute(*this.W, TmplDataAll{
		System: this.TmplGetSystemData(),
	})
}

func (this *Wrapper) printTmplPageError(perr error) {
	tmpl, err := template.New("template").Parse(string(templates.PageTmplError))
	if err != nil {
		(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		(*this.W).WriteHeader(http.StatusInternalServerError)
		(*this.W).Header().Set("Content-Type", "text/html")
		(*this.W).Write([]byte("<h1>Critical engine error!</h1>"))
		(*this.W).Write([]byte("<h2>" + perr.Error() + "</h2>"))
		return
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.W).WriteHeader(http.StatusInternalServerError)
	(*this.W).Header().Set("Content-Type", "text/html")
	tmpl.Execute(*this.W, TmplDataAll{
		System: this.TmplGetSystemData(),
		Data: tmplDataErrorMsg{
			ErrorMessage: perr.Error(),
		},
	})
}

func (this *Wrapper) printEnginePageError(perr error) {
	tmpl, err := template.New("template").Parse(string(templates.PageEngError))
	if err != nil {
		(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		(*this.W).WriteHeader(http.StatusInternalServerError)
		(*this.W).Header().Set("Content-Type", "text/html")
		(*this.W).Write([]byte("<h1>Critical engine error!</h1>"))
		(*this.W).Write([]byte("<h2>" + perr.Error() + "</h2>"))
		return
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.W).WriteHeader(http.StatusInternalServerError)
	(*this.W).Header().Set("Content-Type", "text/html")
	tmpl.Execute(*this.W, TmplDataAll{
		System: this.TmplGetSystemData(),
		Data: tmplDataErrorMsg{
			ErrorMessage: perr.Error(),
		},
	})
}
