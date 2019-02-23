package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
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

func IsDirExists(path string) bool {
	if IsFileExists(path) && IsDir(path) {
		return true
	}
	return false
}

func IsValidEmail(email string) bool {
	regexpe := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regexpe.MatchString(email)
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

func GetMd5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SystemRenderTemplate(w http.ResponseWriter, c []byte, d interface{}) {
	tmpl, err := template.New("template").Parse(string(c))
	if err != nil {
		SystemErrorPageEngine(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, consts.TmplData{
		System: GetTmplSystemData(),
		Data:   d,
	})
}

func SystemErrorPageEngine(w http.ResponseWriter, err error) {
	if tmpl, e := template.New("template").Parse(string(assets.TmplPageErrorEngine)); e == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, consts.TmplData{
			System: GetTmplSystemData(),
			Data:   GetTmplError(err),
		})
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Critical engine error</h1>"))
	w.Write([]byte("<h2>" + err.Error() + "</h2>"))
}

func SystemErrorPage404(w http.ResponseWriter) {
	tmpl, err := template.New("template").Parse(string(assets.TmplPageError404))
	if err != nil {
		SystemErrorPageEngine(w, err)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, consts.TmplData{
		System: GetTmplSystemData(),
		Data:   nil,
	})
}

func UrlToArray(url string) []string {
	url_buff := url
	if len(url_buff) >= 1 && url_buff[:1] == "/" {
		url_buff = url_buff[1:]
	}
	if len(url_buff) >= 1 && url_buff[len(url_buff)-1:] == "/" {
		url_buff = url_buff[:len(url_buff)-1]
	}
	if url_buff == "" {
		return []string{}
	} else {
		return strings.Split(url_buff, "/")
	}
}

func IntToStr(num int) string {
	return fmt.Sprintf("%d", num)
}
