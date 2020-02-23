package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang-fave/engine/consts"
)

func Expect(t *testing.T, actual, expect interface{}) {
	if actual != expect {
		t.Fatalf("\033[0;33mExpected \033[0;32m`(%T) %v`\033[0;33m but got \033[0;31m`(%T) %v`\033[0m",
			expect, expect, actual, actual)
	}
}

func TestIsFileExists(t *testing.T) {
	Expect(t, IsFileExists("./../../support"), true)
	Expect(t, IsFileExists("./../../support/some-file.txt"), true)
	Expect(t, IsFileExists("./../../support/no-existed-file"), false)
}

func TestIsRegularFileExists(t *testing.T) {
	Expect(t, IsRegularFileExists("./../../support/some-file.txt"), true)
	Expect(t, IsRegularFileExists("./../../support/no-existed-file"), false)
	Expect(t, IsRegularFileExists("./../../support"), false)
}

func TestIsDir(t *testing.T) {
	Expect(t, IsDir("./../../support"), true)
	Expect(t, IsDir("./../../support/some-file.txt"), false)
	Expect(t, IsDir("./../../support/no-existed-dir"), false)
}

func TestIsDirExists(t *testing.T) {
	Expect(t, IsDirExists("./../../support"), true)
	Expect(t, IsDirExists("./../../support/some-file.txt"), false)
	Expect(t, IsDirExists("./../../support/no-existed-dir"), false)
}

func TestIsNumeric(t *testing.T) {
	Expect(t, IsNumeric("12345"), true)
	Expect(t, IsNumeric("string"), false)
}

func TestIsFloat(t *testing.T) {
	Expect(t, IsFloat("12345"), true)
	Expect(t, IsFloat("1.23"), true)
	Expect(t, IsFloat("1,23"), false)
	Expect(t, IsFloat("string"), false)
}

func TestIsValidEmail(t *testing.T) {
	Expect(t, IsValidEmail("test@gmail.com"), true)
	Expect(t, IsValidEmail("test@yandex.ru"), true)
	Expect(t, IsValidEmail("test@ya.ru"), true)
	Expect(t, IsValidEmail("test@test"), false)
}

func TestIsValidAlias(t *testing.T) {
	Expect(t, IsValidAlias("/"), true)
	Expect(t, IsValidAlias("/some-page/"), true)
	Expect(t, IsValidAlias("/some-page.html"), true)
	Expect(t, IsValidAlias("/some-page.html/"), true)
	Expect(t, IsValidAlias(""), false)
	Expect(t, IsValidAlias("some-page"), false)
	Expect(t, IsValidAlias("/some page/"), false)

	Expect(t, IsValidAlias("/cp"), false)
	Expect(t, IsValidAlias("/cp/"), false)
	Expect(t, IsValidAlias("/cp/some"), false)
	Expect(t, IsValidAlias("/cp-1"), true)
	Expect(t, IsValidAlias("/cp-some"), true)

	Expect(t, IsValidAlias("/blog"), false)
	Expect(t, IsValidAlias("/blog/"), false)
	Expect(t, IsValidAlias("/blog/some"), false)
	Expect(t, IsValidAlias("/blog-1"), true)
	Expect(t, IsValidAlias("/blog-some"), true)

	Expect(t, IsValidAlias("/shop"), false)
	Expect(t, IsValidAlias("/shop/"), false)
	Expect(t, IsValidAlias("/shop/some"), false)
	Expect(t, IsValidAlias("/shop-1"), true)
	Expect(t, IsValidAlias("/shop-some"), true)

	Expect(t, IsValidAlias("/api"), false)
	Expect(t, IsValidAlias("/api/"), false)
	Expect(t, IsValidAlias("/api/some"), false)
	Expect(t, IsValidAlias("/api-1"), true)
	Expect(t, IsValidAlias("/api-some"), true)
}

func TestIsValidSingleAlias(t *testing.T) {
	Expect(t, IsValidSingleAlias("some-category"), true)
	Expect(t, IsValidSingleAlias("some-category-12345"), true)
	Expect(t, IsValidSingleAlias("some_category_12345"), true)
	Expect(t, IsValidSingleAlias(""), false)
	Expect(t, IsValidSingleAlias("/"), false)
	Expect(t, IsValidSingleAlias("/some-category/"), false)
	Expect(t, IsValidSingleAlias("some-category.html"), false)
	Expect(t, IsValidSingleAlias("some category"), false)
}

func TestFixPath(t *testing.T) {
	Expect(t, FixPath(""), "")
	Expect(t, FixPath("/"), "")
	Expect(t, FixPath("./dir"), "./dir")
	Expect(t, FixPath("./dir/"), "./dir")
	Expect(t, FixPath("\\dir"), "\\dir")
	Expect(t, FixPath("\\dir\\"), "\\dir")
}

func TestExtractHostPort(t *testing.T) {
	h, p := ExtractHostPort("localhost:8080", false)
	Expect(t, h, "localhost")
	Expect(t, p, "8080")

	h, p = ExtractHostPort("localhost:80", false)
	Expect(t, h, "localhost")
	Expect(t, p, "80")

	h, p = ExtractHostPort("localhost", false)
	Expect(t, h, "localhost")
	Expect(t, p, "80")

	h, p = ExtractHostPort("localhost", true)
	Expect(t, h, "localhost")
	Expect(t, p, "443")
}

func TestGetAssetsUrl(t *testing.T) {
	Expect(t, GetAssetsUrl("style.css"), "/style.css?v="+consts.ServerVersion)
}

func TestGetTmplSystemData(t *testing.T) {
	Expect(t, GetTmplSystemData("module", "module"), consts.TmplSystem{
		CpModule:             "module",
		CpSubModule:          "module",
		InfoVersion:          consts.ServerVersion,
		PathCssBootstrap:     "/assets/bootstrap.css?v=" + consts.ServerVersion,
		PathCssCpCodeMirror:  "/assets/cp/tmpl-editor/codemirror.css?v=" + consts.ServerVersion,
		PathCssCpStyles:      "/assets/cp/styles.css?v=" + consts.ServerVersion,
		PathCssCpWysiwygPell: "/assets/cp/wysiwyg/pell.css?v=" + consts.ServerVersion,
		PathCssStyles:        "/assets/sys/styles.css?v=" + consts.ServerVersion,
		PathJsBootstrap:      "/assets/bootstrap.js?v=" + consts.ServerVersion,
		PathJsCpCodeMirror:   "/assets/cp/tmpl-editor/codemirror.js?v=" + consts.ServerVersion,
		PathJsCpScripts:      "/assets/cp/scripts.js?v=" + consts.ServerVersion,
		PathJsCpWysiwygPell:  "/assets/cp/wysiwyg/pell.js?v=" + consts.ServerVersion,
		PathJsJquery:         "/assets/jquery.js?v=" + consts.ServerVersion,
		PathJsPopper:         "/assets/popper.js?v=" + consts.ServerVersion,
		PathSvgLogo:          "/assets/sys/logo.svg?v=" + consts.ServerVersion,
		PathThemeScripts:     "/assets/theme/scripts.js?v=" + consts.ServerVersion,
		PathThemeStyles:      "/assets/theme/styles.css?v=" + consts.ServerVersion,
		PathIcoFav:           "/assets/sys/fave.ico?v=" + consts.ServerVersion,
		PathCssLightGallery:  "/assets/lightgallery.css?v=" + consts.ServerVersion,
		PathJsLightGallery:   "/assets/lightgallery.js?v=" + consts.ServerVersion,
	})
}

func TestGetTmplError(t *testing.T) {
	Expect(t, GetTmplError(errors.New("some error")), consts.TmplError{
		ErrorMessage: "some error",
	})
}

func TestGetMd5(t *testing.T) {
	Expect(t, GetMd5("some string"), "5ac749fbeec93607fc28d666be85e73a")
}

func TestGetCurrentUnixTimestamp(t *testing.T) {
	Expect(t, GetCurrentUnixTimestamp(), int64(time.Now().Unix()))
}

func TestSystemRenderTemplate(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SystemRenderTemplate(w, []byte(`ok`), nil, "module", "module")
	}).ServeHTTP(recorder, request)
	Expect(t, recorder.Code, 200)
	Expect(t, recorder.Body.String(), `ok`)
}

func TestSystemErrorPageEngine(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SystemErrorPageEngine(w, errors.New("Test error"))
	}).ServeHTTP(recorder, request)
	Expect(t, recorder.Code, http.StatusInternalServerError)
	Expect(t, strings.Contains(recorder.Body.String(), "Engine Error"), true)
	Expect(t, strings.Contains(recorder.Body.String(), "Test error"), true)
}

func TestSystemErrorPageTemplate(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SystemErrorPageTemplate(w, errors.New("Test error"))
	}).ServeHTTP(recorder, request)
	Expect(t, recorder.Code, http.StatusInternalServerError)
	Expect(t, strings.Contains(recorder.Body.String(), "Template Error"), true)
	Expect(t, strings.Contains(recorder.Body.String(), "Test error"), true)
}

func TestSystemErrorPage404(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SystemErrorPage404(w)
	}).ServeHTTP(recorder, request)
	Expect(t, recorder.Code, http.StatusNotFound)
	Expect(t, strings.Contains(recorder.Body.String(), "404 Not Found"), true)
}

func TestUrlToArray(t *testing.T) {
	a := UrlToArray("/some/url")
	Expect(t, len(a), 2)
	Expect(t, a[0], "some")
	Expect(t, a[1], "url")

	a = UrlToArray("/some/url/")
	Expect(t, len(a), 2)
	Expect(t, a[0], "some")
	Expect(t, a[1], "url")

	a = UrlToArray("/some/url?a=1&b=2")
	Expect(t, len(a), 2)
	Expect(t, a[0], "some")
	Expect(t, a[1], "url")

	a = UrlToArray("/some/url/?a=1&b=2")
	Expect(t, len(a), 2)
	Expect(t, a[0], "some")
	Expect(t, a[1], "url")
}

func TestIntToStr(t *testing.T) {
	Expect(t, IntToStr(2000), "2000")
}

func TestInt64ToStr(t *testing.T) {
	Expect(t, Int64ToStr(2000), "2000")
}

func TestStrToInt(t *testing.T) {
	Expect(t, StrToInt("2000"), 2000)
	Expect(t, StrToInt("string"), 0)
}

func TestFloat64ToStr(t *testing.T) {
	Expect(t, Float64ToStr(0), "0.00")
	Expect(t, Float64ToStr(0.5), "0.50")
	Expect(t, Float64ToStr(15.8100), "15.81")
}

func TestFloat64ToStrF(t *testing.T) {
	Expect(t, Float64ToStrF(0, "%.4f"), "0.0000")
	Expect(t, Float64ToStrF(0.5, "%.4f"), "0.5000")
	Expect(t, Float64ToStrF(15.8100, "%.4f"), "15.8100")
}

func TestStrToFloat64(t *testing.T) {
	Expect(t, StrToFloat64("0.00"), 0.0)
	Expect(t, StrToFloat64("0.5"), 0.5)
	Expect(t, StrToFloat64("0.50"), 0.5)
	Expect(t, StrToFloat64("15.8100"), 15.81)
	Expect(t, StrToFloat64("15.8155"), 15.8155)
}

func TestGenerateAlias(t *testing.T) {
	Expect(t, GenerateAlias(""), "")
	Expect(t, GenerateAlias("Some page name"), "/some-page-name/")
	Expect(t, GenerateAlias("Some page name 2"), "/some-page-name-2/")
	Expect(t, GenerateAlias("Какая-то страница"), "/kakayato-stranica/")
	Expect(t, GenerateAlias("Какая-то страница 2"), "/kakayato-stranica-2/")
}

func TestGenerateSingleAlias(t *testing.T) {
	Expect(t, GenerateSingleAlias(""), "")
	Expect(t, GenerateSingleAlias("Some category name"), "some-category-name")
	Expect(t, GenerateSingleAlias("Some category name 2"), "some-category-name-2")
	Expect(t, GenerateSingleAlias("Какая-то категория"), "kakayato-kategoriya")
	Expect(t, GenerateSingleAlias("Какая-то категория 2"), "kakayato-kategoriya-2")
}

func TestUnixTimestampToMySqlDateTime(t *testing.T) {
	Expect(t, UnixTimestampToMySqlDateTime(1551741275), "2019-03-05 01:14:35")
}

func TestUnixTimestampToFormat(t *testing.T) {
	Expect(t, UnixTimestampToFormat(1551741275, "2006/01/02 15:04"), "2019/03/05 01:14")
}

func TestExtractGetParams(t *testing.T) {
	Expect(t, ExtractGetParams("/some-url"), "")
	Expect(t, ExtractGetParams("/some-url/"), "")
	Expect(t, ExtractGetParams("/some-url?a=1&b=2"), "?a=1&b=2")
	Expect(t, ExtractGetParams("/some-url/?a=1&b=2"), "?a=1&b=2")
}

func TestJavaScriptVarValue(t *testing.T) {
	Expect(t, JavaScriptVarValue(`It's "string"`), "It&rsquo;s &rdquo;string&rdquo;")
	Expect(t, JavaScriptVarValue(`It is string`), "It is string")
}

func TestInArrayInt(t *testing.T) {
	slice := []int{1, 3, 5, 9, 0}
	Expect(t, InArrayInt(slice, 1), true)
	Expect(t, InArrayInt(slice, 9), true)
	Expect(t, InArrayInt(slice, 2), false)
	Expect(t, InArrayInt(slice, 8), false)
}

func TestInArrayString(t *testing.T) {
	slice := []string{"1", "3", "5", "9", "0"}
	Expect(t, InArrayString(slice, "1"), true)
	Expect(t, InArrayString(slice, "9"), true)
	Expect(t, InArrayString(slice, "2"), false)
	Expect(t, InArrayString(slice, "8"), false)
}

func TestGetPostArrayInt(t *testing.T) {
	request, err := http.NewRequest("POST", "/", strings.NewReader("cats[]=1&cats[]=3&cats[]=5"))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	request.ParseForm()
	arr := GetPostArrayInt("cats[]", request)
	Expect(t, fmt.Sprintf("%T%v", arr, arr), "[]int[1 3 5]")
}

func TestGetPostArrayString(t *testing.T) {
	request, err := http.NewRequest("POST", "/", strings.NewReader("cats[]=1&cats[]=3&cats[]=5"))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	request.ParseForm()
	arr := GetPostArrayString("cats[]", request)
	Expect(t, fmt.Sprintf("%T%v", arr, arr), "[]string[1 3 5]")
}

func TestArrayOfIntToArrayOfString(t *testing.T) {
	res := ArrayOfIntToArrayOfString([]int{1, 3, 5})
	Expect(t, len(res), 3)
	Expect(t, res[0], "1")
	Expect(t, res[1], "3")
	Expect(t, res[2], "5")
}

func TestArrayOfStringToArrayOfInt(t *testing.T) {
	res := ArrayOfStringToArrayOfInt([]string{"1", "3", "5", "abc"})
	Expect(t, len(res), 3)
	Expect(t, res[0], 1)
	Expect(t, res[1], 3)
	Expect(t, res[2], 5)
}

func TestGetImagePlaceholderSrc(t *testing.T) {
	Expect(t, GetImagePlaceholderSrc(), "/assets/sys/placeholder.png")
}

func TestFormatProductPrice(t *testing.T) {
	Expect(t, FormatProductPrice(123.4567, 0, 0), "123")
	Expect(t, FormatProductPrice(123.4567, 0, 1), "124")
	Expect(t, FormatProductPrice(123.4567, 0, 2), "123")

	Expect(t, FormatProductPrice(123.4567, 1, 0), "123.5")
	Expect(t, FormatProductPrice(123.4567, 2, 0), "123.46")
	Expect(t, FormatProductPrice(123.4567, 3, 0), "123.457")
	Expect(t, FormatProductPrice(123.4567, 4, 0), "123.4567")

	Expect(t, FormatProductPrice(123.4567, 1, 1), "124.0")
	Expect(t, FormatProductPrice(123.4567, 2, 1), "124.00")
	Expect(t, FormatProductPrice(123.4567, 3, 1), "124.000")
	Expect(t, FormatProductPrice(123.4567, 4, 1), "124.0000")

	Expect(t, FormatProductPrice(123.4567, 1, 2), "123.0")
	Expect(t, FormatProductPrice(123.4567, 2, 2), "123.00")
	Expect(t, FormatProductPrice(123.4567, 3, 2), "123.000")
	Expect(t, FormatProductPrice(123.4567, 4, 2), "123.0000")
}

func TestSafeFilePath(t *testing.T) {
	Expect(t, SafeFilePath("/test/file"), "/test/file")
	Expect(t, SafeFilePath("/test/../file"), "/test/file")
	Expect(t, SafeFilePath("../test/file"), "/test/file")
	Expect(t, SafeFilePath("/test/file/.."), "/test/file/")
	Expect(t, SafeFilePath("/test/file/./"), "/test/file/")
	Expect(t, SafeFilePath("/test/./file"), "/test/file")
}

func TestIsValidTemplateFileName(t *testing.T) {
	Expect(t, IsValidTemplateFileName("test-template"), true)
	Expect(t, IsValidTemplateFileName("test-123-TEST"), true)
	Expect(t, IsValidTemplateFileName("TEST-123-TEST"), true)
	Expect(t, IsValidTemplateFileName("test template"), false)
	Expect(t, IsValidTemplateFileName("test_template"), false)
	Expect(t, IsValidTemplateFileName("test-template.html"), false)
	Expect(t, IsValidTemplateFileName("test-template.css"), false)
	Expect(t, IsValidTemplateFileName("test@template"), false)
}
