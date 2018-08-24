package wrapper

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	Images "golang-fave/engine/wrapper/resources/images"
	Others "golang-fave/engine/wrapper/resources/others"
	Styles "golang-fave/engine/wrapper/resources/styles"
	Templates "golang-fave/engine/wrapper/resources/templates"
)

func (e *Wrapper) isFileExists(file string) bool {
	if fi, err := os.Stat(file); !os.IsNotExist(err) {
		if err == nil {
			fmode := fi.Mode()
			if !fmode.IsDir() {
				return true
			}
		}
	}
	return false
}

func (e *Wrapper) getFileContentType(file string, of *os.File) string {
	fContentType := ""
	if strings.HasSuffix(file, ".htm") {
		fContentType = "text/html"
	} else if strings.HasSuffix(file, ".html") {
		fContentType = "text/html"
	} else if strings.HasSuffix(file, ".txt") {
		fContentType = "text/plain"
	} else if strings.HasSuffix(file, ".php") {
		fContentType = "text/plain"
	} else if strings.HasSuffix(file, ".css") {
		fContentType = "text/css"
	} else if strings.HasSuffix(file, ".png") {
		fContentType = "image/png"
	} else if strings.HasSuffix(file, ".jpg") {
		fContentType = "image/jpeg"
	} else if strings.HasSuffix(file, ".jpeg") {
		fContentType = "image/jpeg"
	} else {
		fContentType = "application/octet-stream"
	}
	return fContentType
}

func (e *Wrapper) staticResource() bool {
	if e.R.URL.Path == "/assets/sys/styles.css" {
		(*e.W).Header().Set("Content-Type", "text/css")
		(*e.W).Write(Styles.File_assets_sys_styles_css)
		return true
	} else if e.R.URL.Path == "/assets/sys/logo.svg" {
		(*e.W).Header().Set("Content-Type", "image/svg+xml")
		(*e.W).Write(Others.File_assets_sys_logo_svg)
		return true
	} else if e.R.URL.Path == "/assets/sys/bg.png" {
		(*e.W).Header().Set("Content-Type", "image/png")
		(*e.W).Write(Images.File_assets_sys_bg_png)
		return true
	} else if e.R.URL.Path == "/assets/sys/logo.png" {
		(*e.W).Header().Set("Content-Type", "image/png")
		(*e.W).Write(Images.File_assets_sys_logo_png)
		return true
	} else if e.R.URL.Path == "/assets/sys/fave.ico" {
		(*e.W).Header().Set("Content-Type", "image/x-icon")
		(*e.W).Write(Others.File_assets_sys_fave_ico)
		return true
	}
	return false
}

func (e *Wrapper) staticFile() bool {
	file := e.R.URL.Path
	if file == "/" {
		if e.isFileExists(e.DirVhostHome + "/htdocs" + "/index.htm") {
			file = "/index.htm"
		} else if e.isFileExists(e.DirVhostHome + "/htdocs" + "/index.html") {
			file = "/index.html"
		}
	}
	if file != "/" {
		if e.isFileExists(e.DirVhostHome + "/htdocs" + file) {
			of, err := os.Open(e.DirVhostHome + "/htdocs" + file)
			defer of.Close()
			if err == nil {
				fstat, _ := of.Stat()
				fsize := strconv.FormatInt(fstat.Size(), 10)
				(*e.W).Header().Add("Content-Type", e.getFileContentType(file, of))
				(*e.W).Header().Add("Content-Length", fsize)
				of.Seek(0, 0)
				io.Copy(*e.W, of)
				return true
			}
		}
	}
	return false
}

func (e *Wrapper) printPageDefault() {
	(*e.W).Header().Set("Content-Type", "text/html")
	(*e.W).Write(Templates.PageDefault)
}

func (e *Wrapper) printPage404() {
	(*e.W).WriteHeader(http.StatusNotFound)
	(*e.W).Header().Set("Content-Type", "text/html")
	(*e.W).Write(Templates.PageError404)
}
