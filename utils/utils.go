package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

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

func IsNumeric(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

func IsValidEmail(email string) bool {
	regexpe := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regexpe.MatchString(email)
}

func IsValidAlias(alias string) bool {
	regexpeSlash := regexp.MustCompile(`[\/]{2,}`)
	regexpeChars := regexp.MustCompile(`^\/([a-zA-Z0-9\/\-_\.]+)\/?$`)
	return (!regexpeSlash.MatchString(alias) && regexpeChars.MatchString(alias)) || alias == "/"
}

func FixPath(path string) string {
	newPath := strings.TrimSpace(path)
	if len(newPath) <= 0 {
		return newPath
	}
	if newPath[len(newPath)-1] == '/' || newPath[len(newPath)-1] == '\\' {
		newPath = newPath[0 : len(newPath)-1]
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
		PathThemeStyles:  "/assets/theme/styles.css",
		PathThemeScripts: "/assets/theme/scripts.js",
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

func GetCurrentUnixTimestamp() int64 {
	return int64(time.Now().Unix())
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

func SystemErrorPageTemplate(w http.ResponseWriter, err error) {
	if tmpl, e := template.New("template").Parse(string(assets.TmplPageErrorTmpl)); e == nil {
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

	// Remove GET parameters
	i := strings.Index(url_buff, "?")
	if i > -1 {
		url_buff = url_buff[:i]
	}

	// Cut slashes
	if len(url_buff) >= 1 && url_buff[:1] == "/" {
		url_buff = url_buff[1:]
	}
	if len(url_buff) >= 1 && url_buff[len(url_buff)-1:] == "/" {
		url_buff = url_buff[:len(url_buff)-1]
	}

	// Explode
	if url_buff == "" {
		return []string{}
	} else {
		return strings.Split(url_buff, "/")
	}
}

func IntToStr(num int) string {
	return fmt.Sprintf("%d", num)
}

func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err == nil {
		return num
	}
	return 0
}

func GenerateAlias(str string) string {
	if str == "" {
		return ""
	}

	strc := utf16.Encode([]rune(str))

	lat := []string{"EH", "I", "i", "#", "eh", "A", "B", "V", "G", "D", "E", "JO", "ZH", "Z", "I", "JJ", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "F", "KH", "C", "CH", "SH", "SHH", "'", "Y", "", "EH", "YU", "YA", "a", "b", "v", "g", "d", "e", "jo", "zh", "z", "i", "jj", "k", "l", "m", "n", "o", "p", "r", "s", "t", "u", "f", "kh", "c", "ch", "sh", "shh", "", "y", "", "eh", "yu", "ya", "", "", "-", "-", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "[", "]", "a", "s", "d", "f", "g", "h", "j", "k", "l", ";", "'", "z", "x", "c", "v", "b", "n", "m", ",", ".", "/", "-", "-", ":", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "A", "S", "D", "F", "G", "H", "J", "K", "L", "Z", "X", "C", "V", "B", "N", "M"}
	cyr := []string{"Є", "І", "і", "№", "є", "А", "Б", "В", "Г", "Д", "Е", "Ё", "Ж", "З", "И", "Й", "К", "Л", "М", "Н", "О", "П", "Р", "С", "Т", "У", "Ф", "Х", "Ц", "Ч", "Ш", "Щ", "Ъ", "Ы", "Ь", "Э", "Ю", "Я", "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я", "«", "»", "—", " ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "", "", "a", "s", "d", "f", "g", "h", "j", "k", "l", "", "", "z", "x", "c", "v", "b", "n", "m", "", "", "", "(", ")", "", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "A", "S", "D", "F", "G", "H", "J", "K", "L", "Z", "X", "C", "V", "B", "N", "M"}

	var alias string = ""
	for i := 0; i < len(strc); i++ {
		for j := 0; j < len(cyr); j++ {
			if string(strc[i]) == cyr[j] {
				alias += lat[j]
			}
		}
	}
	alias = strings.ToLower(alias)

	// Cut repeated chars "-"
	if reg, err := regexp.Compile("[\\-]+"); err == nil {
		alias = strings.Trim(reg.ReplaceAllString(alias, "-"), "-")
	}

	alias = "/" + alias + "/"

	// Cut repeated chars "/"
	if reg, err := regexp.Compile("[/]+"); err == nil {
		alias = reg.ReplaceAllString(alias, "/")
	}

	return alias
}

func UnixTimestampToMySqlDateTime(sec int64) string {
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

func UnixTimestampToFormat(sec int64, format string) string {
	return time.Unix(sec, 0).Format(format)
}

func ExtractGetParams(str string) string {
	i := strings.Index(str, "?")
	if i == -1 {
		return ""
	}
	return "?" + str[i+1:]
}

func JavaScriptVarValue(str string) string {
	return strings.Replace(
		strings.Replace(str, `'`, `&rsquo;`, -1),
		`"`,
		`&rdquo;`,
		-1,
	)
}
