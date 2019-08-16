package consts

import (
	"html/template"
)

const AssetsPath = "assets"
const AssetsVersion = "40"
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

// Wysiwyg editor
const AssetsCpWysiwygPellCss = AssetsPath + "/cp/wysiwyg/pell.css"
const AssetsCpWysiwygPellJs = AssetsPath + "/cp/wysiwyg/pell.js"

// CodeMirror template editor
const AssetsCpCodeMirrorCss = AssetsPath + "/cp/tmpl-editor/codemirror.css"
const AssetsCpCodeMirrorJs = AssetsPath + "/cp/tmpl-editor/codemirror.js"

// Make global for other packages
var ParamDebug bool
var ParamHost string
var ParamKeepAlive bool
var ParamPort int
var ParamWwwDir string

// For admin panel
type BreadCrumb struct {
	Name string
	Link string
}

// Template data
type TmplSystem struct {
	CpModule             string
	CpSubModule          string
	InfoVersion          string
	PathCssBootstrap     string
	PathCssCpCodeMirror  string
	PathCssCpStyles      string
	PathCssCpWysiwygPell string
	PathCssStyles        string
	PathJsBootstrap      string
	PathJsCpCodeMirror   string
	PathJsCpScripts      string
	PathJsCpWysiwygPell  string
	PathJsJquery         string
	PathJsPopper         string
	PathSvgLogo          string
	PathThemeScripts     string
	PathThemeStyles      string
	PathIcoFav           string
}

type TmplError struct {
	ErrorMessage string
}

type TmplData struct {
	System TmplSystem
	Data   interface{}
}

type TmplDataCpBase struct {
	BodyClasses        string
	Caption            string
	Content            template.HTML
	ModuleCurrentAlias string
	NavBarModules      template.HTML
	NavBarModulesSys   template.HTML
	SidebarLeft        template.HTML
	SidebarRight       template.HTML
	UserAvatarLink     string
	UserEmail          string
	UserFirstName      string
	UserId             int
	UserLastName       string
	UserPassword       string
	Title              string
}
