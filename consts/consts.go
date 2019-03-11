package consts

import (
	"html/template"
)

const ServerVersion = "1.0.2"
const AssetsVersion = "16"
const AssetsPath = "assets"
const DirIndexFile = "index.html"

// Bootstrap resources
const AssetsBootstrapCss = AssetsPath + "/bootstrap.css"
const AssetsBootstrapJs = AssetsPath + "/bootstrap.js"
const AssetsJqueryJs = AssetsPath + "/jquery.js"
const AssetsPopperJs = AssetsPath + "/popper.js"

// System resources
const AssetsCpScriptsJs = AssetsPath + "/cp/scripts.js"
const AssetsCpStylesCss = AssetsPath + "/cp/styles.css"
const AssetsSysBgPng = AssetsPath + "/sys/bg.png"
const AssetsSysFaveIco = AssetsPath + "/sys/fave.ico"
const AssetsSysLogoPng = AssetsPath + "/sys/logo.png"
const AssetsSysLogoSvg = AssetsPath + "/sys/logo.svg"
const AssetsSysStylesCss = AssetsPath + "/sys/styles.css"

// Make global for other packages
var ParamHost string
var ParamPort int
var ParamWwwDir string
var ParamDebug bool
var ParamKeepAlive bool

// For admin panel
type BreadCrumb struct {
	Name string
	Link string
}

// Template data
type TmplSystem struct {
	PathIcoFav       string
	PathSvgLogo      string
	PathCssStyles    string
	PathCssCpStyles  string
	PathCssBootstrap string
	PathJsJquery     string
	PathJsPopper     string
	PathJsBootstrap  string
	PathJsCpScripts  string
	PathThemeStyles  string
	PathThemeScripts string
	InfoVersion      string
}

type TmplError struct {
	ErrorMessage string
}

type TmplData struct {
	System TmplSystem
	Data   interface{}
}

type TmplDataCpBase struct {
	Title              string
	Caption            string
	BodyClasses        string
	UserId             int
	UserFirstName      string
	UserLastName       string
	UserEmail          string
	UserPassword       string
	UserAvatarLink     string
	NavBarModules      template.HTML
	NavBarModulesSys   template.HTML
	ModuleCurrentAlias string
	SidebarLeft        template.HTML
	Content            template.HTML
	SidebarRight       template.HTML
}

type TmplDataMainMenuItem struct {
	Name   string
	Link   string
	Active bool
}

type TmplDataModIndex struct {
	Name            string
	Alias           string
	Content         template.HTML
	MetaTitle       string
	MetaKeywords    string
	MetaDescription string
}
