package wrapper

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang-fave/constants"
	"golang-fave/engine/sessions"
	templates "golang-fave/engine/wrapper/resources/templates"
)

type handleRun func(wrapper *Wrapper) bool

type tmplDataSystem struct {
	PathIcoFav       string
	PathSvgLogo      string
	PathCssStyles    string
	PathCssCpStyles  string
	PathCssBootstrap string
	PathJsJquery     string
	PathJsPopper     string
	PathJsBootstrap  string
	PathJsCpScripts  string
	BlockModalSysMsg template.HTML
}

type tmplDataAll struct {
	System tmplDataSystem
	Data   interface{}
}

type Wrapper struct {
	W            *http.ResponseWriter
	R            *http.Request
	VHost        string
	Port         string
	DirWww       string
	DirVHostHome string
	RemoteIp     string
	LoggerAcc    *log.Logger
	LoggerErr    *log.Logger
	Session      *sessions.Session
}

func (this *Wrapper) tmplGetSystemData() tmplDataSystem {
	return tmplDataSystem{
		PathIcoFav:       this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/fave.ico?v=" + constants.AssetsVersion,
		PathSvgLogo:      this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/logo.svg?v=" + constants.AssetsVersion,
		PathCssStyles:    this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/styles.css?v=" + constants.AssetsVersion,
		PathCssCpStyles:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/cp/styles.css?v=" + constants.AssetsVersion,
		PathCssBootstrap: this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/bootstrap.css?v=" + constants.AssetsVersion,
		PathJsJquery:     this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/jquery.js?v=" + constants.AssetsVersion,
		PathJsPopper:     this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/popper.js?v=" + constants.AssetsVersion,
		PathJsBootstrap:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/bootstrap.js?v=" + constants.AssetsVersion,
		PathJsCpScripts:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/cp/scripts.js?v=" + constants.AssetsVersion,
		BlockModalSysMsg: template.HTML(templates.BlockModalSysMsg),
	}
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, port string, wwwdir string, vhosthome string) *Wrapper {
	return &Wrapper{
		VHost:        vhost,
		Port:         port,
		DirWww:       wwwdir,
		DirVHostHome: vhosthome,
		W:            w,
		R:            r,
	}
}

func (this *Wrapper) Run(hRun handleRun) {
	// Populate some values
	this.RemoteIp = this.R.RemoteAddr

	// Create loggers
	this.LoggerAcc = log.New(os.Stdout, this.VHost+", ", log.LstdFlags)
	this.LoggerErr = log.New(os.Stdout, this.VHost+", ", log.LstdFlags)

	// Attach file for access log
	if !constants.Debug {
		acclogfile, acclogfileerr := os.OpenFile(this.DirVHostHome+"/logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if acclogfileerr == nil {
			defer acclogfile.Close()
			this.LoggerAcc.SetOutput(acclogfile)
		}
	}

	// Attach file for access log
	if !constants.Debug {
		errlogfile, errlogfileerr := os.OpenFile(this.DirVHostHome+"/logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if errlogfileerr == nil {
			defer errlogfile.Close()
			this.LoggerErr.SetOutput(errlogfile)
		}
	}

	// Fix remote IP
	if strings.ContainsRune(this.R.RemoteAddr, ':') {
		this.RemoteIp, _, _ = net.SplitHostPort(this.R.RemoteAddr)
	}

	// Redirect to main domain
	if this.redirectToMainDomain() {
		this.Log("301")
		return
	}

	// Static resource
	if this.staticResource() {
		this.Log("200")
		return
	}

	// Static file
	if this.staticFile() {
		this.Log("200")
		return
	}

	// Friendly search engine url
	if this.redirectSeoFix() {
		this.Log("301")
		return
	}

	// Create and load session
	this.Session = sessions.New(this.W, this.R, this.VHost, this.DirVHostHome, this.RemoteIp)
	this.Session.Load()

	// Set session vars
	if !this.Session.IsSetInt("UserId") {
		this.Session.SetInt("UserId", 0)
	}
	if !this.Session.IsSetBool("IsLogged") {
		this.Session.SetBool("IsLogged", false)
	}

	// Logic
	if hRun != nil {
		if hRun(this) {
			this.Log("200")
			this.Session.Save()
			return
		}
	}

	// Show default page
	if this.R.URL.Path == "/" {
		this.Log("200")
		this.printPageDefault()
	} else {
		this.LogError("404")
		this.printPage404()
	}
}

func (this *Wrapper) Log(value string) {
	this.LoggerAcc.Println("[ACC] [" + this.R.Method + "] [" + value + "] [" + this.RemoteIp +
		"] [" + this.R.URL.Scheme + "://" + this.R.Host + this.R.URL.RequestURI() +
		"] [" + this.R.Header.Get("User-Agent") + "]")
}

func (this *Wrapper) LogError(value string) {
	this.LoggerErr.Println("[ERR] [" + this.R.Method + "] [" + value + "] [" + this.RemoteIp +
		"] [" + this.R.URL.Scheme + "://" + this.R.Host + this.R.URL.RequestURI() +
		"] [" + this.R.Header.Get("User-Agent") + "]")
}

func (this *Wrapper) TmplFrontEnd(tname string, data interface{}) bool {
	tmpl, err := template.ParseFiles(
		this.DirVHostHome+"/template"+"/"+tname+".html",
		this.DirVHostHome+"/template"+"/header.html",
		this.DirVHostHome+"/template"+"/sidebar.html",
		this.DirVHostHome+"/template"+"/footer.html",
	)
	if err != nil {
		this.printTmplPageError(err)
		return true
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	tmpl.Execute(*this.W, tmplDataAll{
		System: this.tmplGetSystemData(),
		Data:   data,
	})
	return true
}

func (this *Wrapper) TmplBackEnd(tcont []byte, data interface{}) bool {
	tmpl, err := template.New("template").Parse(string(tcont))
	if err != nil {
		this.printTmplPageError(err)
		return true
	}
	(*this.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	tmpl.Execute(*this.W, tmplDataAll{
		System: this.tmplGetSystemData(),
		Data:   data,
	})
	return true
}

func (this *Wrapper) EngineErrMsgOnError(err error) bool {
	if err != nil {
		this.printEnginePageError(err)
		return true
	}
	return false
}
