package wrapper

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang-fave/engine/sessions"
	templates "golang-fave/engine/wrapper/resources/templates"
)

const C_AssetsVersion = "3"

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
	DirVhostHome string
	RemoteIp     string
	LoggerAcc    *log.Logger
	LoggerErr    *log.Logger
	Session      *sessions.Session
	Debug        bool
}

func (this *Wrapper) tmplGetSystemData() tmplDataSystem {
	version := "?v=" + C_AssetsVersion
	return tmplDataSystem{
		PathIcoFav:       this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/fave.ico" + version,
		PathSvgLogo:      this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/logo.svg" + version,
		PathCssStyles:    this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/styles.css" + version,
		PathCssCpStyles:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/cp/styles.css" + version,
		PathCssBootstrap: this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/bootstrap.css" + version,
		PathJsJquery:     this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/jquery.js" + version,
		PathJsPopper:     this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/popper.js" + version,
		PathJsBootstrap:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/sys/bootstrap.js" + version,
		PathJsCpScripts:  this.R.URL.Scheme + "://" + this.R.Host + "/assets/cp/scripts.js" + version,
		BlockModalSysMsg: template.HTML(templates.BlockModalSysMsg),
	}
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, port string, wwwdir string, vhosthome string, debug bool) *Wrapper {
	return &Wrapper{
		VHost:        vhost,
		Port:         port,
		DirWww:       wwwdir,
		DirVhostHome: vhosthome,
		W:            w,
		R:            r,
		Debug:        debug,
	}
}

func (this *Wrapper) Run(hRun handleRun) {
	// Populate some values
	this.RemoteIp = this.R.RemoteAddr

	// Create loggers
	this.LoggerAcc = log.New(os.Stdout, this.VHost+", ", log.LstdFlags)
	this.LoggerErr = log.New(os.Stdout, this.VHost+", ", log.LstdFlags)

	// Attach file for access log
	if !this.Debug {
		acclogfile, acclogfileerr := os.OpenFile(this.DirVhostHome+"/logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if acclogfileerr == nil {
			defer acclogfile.Close()
			this.LoggerAcc.SetOutput(acclogfile)
		}
	}

	// Attach file for access log
	if !this.Debug {
		errlogfile, errlogfileerr := os.OpenFile(this.DirVhostHome+"/logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	this.Session = sessions.New(this.W, this.R, this.VHost, this.DirVhostHome, this.RemoteIp)
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
		this.DirVhostHome+"/template"+"/"+tname+".html",
		this.DirVhostHome+"/template"+"/header.html",
		this.DirVhostHome+"/template"+"/sidebar.html",
		this.DirVhostHome+"/template"+"/footer.html",
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
