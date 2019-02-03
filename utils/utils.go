package utils

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"golang-fave/assets"
	"golang-fave/consts"
)

func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err == nil {
			return true
		}
	}
	return false
}

func IsDir(filename string) bool {
	if st, err := os.Stat(filename); !os.IsNotExist(err) {
		if err == nil {
			if st.Mode().IsDir() {
				return true
			}
		}
	}
	return false
}

func IsHostDirExists(path string) bool {
	if IsFileExists(path) && IsDir(path) {
		return true
	}
	return false
}

func FixPath(path string) string {
	newPath := strings.TrimSpace(path)
	if len(newPath) <= 0 {
		return newPath
	}
	if newPath[len(newPath)-1] == '/' || newPath[len(newPath)-1] == '\\' {
		newPath = newPath[0 : len(newPath)-2]
	}
	return newPath
}

func ExtractHostPort(host string, https bool) (string, string) {
	h := host
	p := "80"
	if https {
		p = "443"
	}
	i := strings.Index(h, ":")
	if i > -1 {
		p = h[i+1:]
		h = h[0:i]
	}
	return h, p
}

func GetAssetsUrl(filename string) string {
	return "/" + filename + "?v=" + consts.AssetsVersion
}

func GetTmplSystemData() consts.TmplSystem {
	return consts.TmplSystem{
		PathIcoFav:       GetAssetsUrl(consts.AssetsSysFaveIco),
		PathSvgLogo:      GetAssetsUrl(consts.AssetsSysLogoSvg),
		PathCssStyles:    GetAssetsUrl(consts.AssetsSysStylesCss),
		PathCssCpStyles:  GetAssetsUrl(consts.AssetsCpStylesCss),
		PathCssBootstrap: GetAssetsUrl(consts.AssetsBootstrapCss),
		PathJsJquery:     GetAssetsUrl(consts.AssetsJqueryJs),
		PathJsPopper:     GetAssetsUrl(consts.AssetsPopperJs),
		PathJsBootstrap:  GetAssetsUrl(consts.AssetsBootstrapJs),
		PathJsCpScripts:  GetAssetsUrl(consts.AssetsCpScriptsJs),
	}
}

func GetTmplError(err error) consts.TmplError {
	return consts.TmplError{
		ErrorMessage: err.Error(),
	}
}

func SystemErrorPage(w http.ResponseWriter, err error) {
	if tmpl, errr := template.New("template").Parse(string(assets.TmplPageErrorEngine)); errr == nil {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, consts.TmplData{
			System: GetTmplSystemData(),
			Data:   GetTmplError(err),
		})
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Critical engine error</h1>"))
	w.Write([]byte("<h2>" + err.Error() + "</h2>"))
}
