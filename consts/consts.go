package consts

import (
	"html/template"
)

const Debug = true
const ServerVersion = "1.0.0"
const AssetsVersion = "1"
const AssetsPath = "assets"
const DirIndexFile = "index.html"

// Bootstrap resources
const AssetsBootstrapCss = AssetsPath + "/bootstrap.css"
const AssetsBootstrapJs = AssetsPath + "/bootstrap.js"
const AssetsJqueryJs = AssetsPath + "/jquery.js"
const AssetsPopperJs = AssetsPath + "/popper.js"

// System resources
const AssetsSysBg404Jpg = AssetsPath + "/sys/bg404.jpg"
const AssetsCpScriptsJs = AssetsPath + "/cp/scripts.js"
const AssetsCpStylesCss = AssetsPath + "/cp/styles.css"
const AssetsSysBgPng = AssetsPath + "/sys/bg.png"
const AssetsSysFaveIco = AssetsPath + "/sys/fave.ico"
const AssetsSysLogoPng = AssetsPath + "/sys/logo.png"
const AssetsSysLogoSvg = AssetsPath + "/sys/logo.svg"
const AssetsSysStylesCss = AssetsPath + "/sys/styles.css"

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
	MetaTitle       string
	MetaKeywords    string
	MetaDescription string
	MainMenuItems   []TmplDataMainMenuItem
}
