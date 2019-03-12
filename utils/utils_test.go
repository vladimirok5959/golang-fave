package utils

import (
	"errors"
	"testing"
	"time"

	"golang-fave/consts"
)

func Expect(t *testing.T, actual, expect interface{}) {
	if actual != expect {
		t.Fatalf("\033[0;33mExpected \033[0;32m`%v`\033[0;33m but got \033[0;31m`%v`\033[0m",
			expect, actual)
	}
}

func TestIsFileExists(t *testing.T) {
	Expect(t, IsFileExists("./../testdata/some-file.txt"), true)
	Expect(t, IsFileExists("./../testdata/no-existed-file"), false)
}

func TestIsDir(t *testing.T) {
	Expect(t, IsDir("./../testdata"), true)
	Expect(t, IsDir("./../testdata/some-file.txt"), false)
	Expect(t, IsDir("./../testdata/no-existed-dir"), false)
}

func TestIsDirExists(t *testing.T) {
	Expect(t, IsDirExists("./../testdata"), true)
	Expect(t, IsDirExists("./../testdata/some-file.txt"), false)
	Expect(t, IsDirExists("./../testdata/no-existed-dir"), false)
}

func TestIsNumeric(t *testing.T) {
	Expect(t, IsNumeric("12345"), true)
	Expect(t, IsNumeric("string"), false)
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
	Expect(t, GetAssetsUrl("style.css"), "/style.css?v="+consts.AssetsVersion)
}

func TestGetTmplSystemData(t *testing.T) {
	Expect(t, GetTmplSystemData(), consts.TmplSystem{
		PathIcoFav:       "/assets/sys/fave.ico?v=" + consts.AssetsVersion,
		PathSvgLogo:      "/assets/sys/logo.svg?v=" + consts.AssetsVersion,
		PathCssStyles:    "/assets/sys/styles.css?v=" + consts.AssetsVersion,
		PathCssCpStyles:  "/assets/cp/styles.css?v=" + consts.AssetsVersion,
		PathCssBootstrap: "/assets/bootstrap.css?v=" + consts.AssetsVersion,
		PathJsJquery:     "/assets/jquery.js?v=" + consts.AssetsVersion,
		PathJsPopper:     "/assets/popper.js?v=" + consts.AssetsVersion,
		PathJsBootstrap:  "/assets/bootstrap.js?v=" + consts.AssetsVersion,
		PathJsCpScripts:  "/assets/cp/scripts.js?v=" + consts.AssetsVersion,
		PathThemeStyles:  "/assets/theme/styles.css",
		PathThemeScripts: "/assets/theme/scripts.js",
		InfoVersion:      consts.ServerVersion,
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
	//
}

func TestSystemErrorPageEngine(t *testing.T) {
	//
}

func TestSystemErrorPageTemplate(t *testing.T) {
	//
}

func TestSystemErrorPage404(t *testing.T) {
	//
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

func TestStrToInt(t *testing.T) {
	Expect(t, StrToInt("2000"), 2000)
	Expect(t, StrToInt("string"), 0)
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
