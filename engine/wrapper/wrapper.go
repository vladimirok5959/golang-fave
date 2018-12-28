package wrapper

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang-fave/engine/actions"
	"golang-fave/engine/sessions"
)

const C_AssetsVersion = "2"

type handleRun func(e *Wrapper) bool

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
	Action       *actions.Action
	Session      *sessions.Session
	Debug        bool
}

func (e *Wrapper) tmplGetSystemData() tmplDataSystem {
	version := "?v=" + C_AssetsVersion
	return tmplDataSystem{
		PathIcoFav:       e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/fave.ico" + version,
		PathSvgLogo:      e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/logo.svg" + version,
		PathCssStyles:    e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/styles.css" + version,
		PathCssCpStyles:  e.R.URL.Scheme + "://" + e.R.Host + "/assets/cp/styles.css" + version,
		PathCssBootstrap: e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/bootstrap.css" + version,
		PathJsJquery:     e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/jquery.js" + version,
		PathJsPopper:     e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/popper.js" + version,
		PathJsBootstrap:  e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/bootstrap.js" + version,
		PathJsCpScripts:  e.R.URL.Scheme + "://" + e.R.Host + "/assets/cp/scripts.js" + version,
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

func (e *Wrapper) Run(hRun handleRun) {
	// Populate some values
	e.RemoteIp = e.R.RemoteAddr

	// Create loggers
	e.LoggerAcc = log.New(os.Stdout, e.VHost+", ", log.LstdFlags)
	e.LoggerErr = log.New(os.Stdout, e.VHost+", ", log.LstdFlags)

	// Attach file for access log
	if !e.Debug {
		acclogfile, acclogfileerr := os.OpenFile(e.DirVhostHome+"/logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if acclogfileerr == nil {
			defer acclogfile.Close()
			e.LoggerAcc.SetOutput(acclogfile)
		}
	}

	// Attach file for access log
	if !e.Debug {
		errlogfile, errlogfileerr := os.OpenFile(e.DirVhostHome+"/logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if errlogfileerr == nil {
			defer errlogfile.Close()
			e.LoggerErr.SetOutput(errlogfile)
		}
	}

	// Fix remote IP
	if strings.ContainsRune(e.R.RemoteAddr, ':') {
		e.RemoteIp, _, _ = net.SplitHostPort(e.R.RemoteAddr)
	}

	// Redirect to main domain
	if e.redirectToMainDomain() {
		e.Log("301")
		return
	}

	// Static resource
	if e.staticResource() {
		e.Log("200")
		return
	}

	// Static file
	if e.staticFile() {
		e.Log("200")
		return
	}

	// Friendly search engine url
	if e.redirectSeoFix() {
		e.Log("301")
		return
	}

	// Create and load session
	e.Session = sessions.New(e.W, e.R, e.VHost, e.DirVhostHome, e.RemoteIp)
	e.Session.Load()

	// Set session vars
	if !e.Session.IsSetInt("UserId") {
		e.Session.SetInt("UserId", 0)
	}
	if !e.Session.IsSetBool("IsLogged") {
		e.Session.SetBool("IsLogged", false)
	}

	// Create action
	e.Action = actions.New(e.W, e.R, e.VHost, e.DirVhostHome, e.RemoteIp)

	// Call action
	if e.Action.Call() {
		e.Session.Save()
		return
	}

	// Logic
	if hRun != nil {
		if hRun(e) {
			e.Session.Save()
			return
		}
	}

	// Show default page
	if e.R.URL.Path == "/" {
		e.Log("200")
		e.printPageDefault()
	} else {
		e.LogError("404")
		e.printPage404()
	}
}

func (e *Wrapper) Log(value string) {
	e.LoggerAcc.Println("[ACC] [" + e.R.Method + "] [" + value + "] [" + e.RemoteIp +
		"] [" + e.R.URL.Scheme + "://" + e.R.Host + e.R.URL.RequestURI() +
		"] [" + e.R.Header.Get("User-Agent") + "]")
}

func (e *Wrapper) LogError(value string) {
	e.LoggerErr.Println("[ERR] [" + e.R.Method + "] [" + value + "] [" + e.RemoteIp +
		"] [" + e.R.URL.Scheme + "://" + e.R.Host + e.R.URL.RequestURI() +
		"] [" + e.R.Header.Get("User-Agent") + "]")
}

func (e *Wrapper) TmplFrontEnd(tname string, data interface{}) bool {
	tmpl, err := template.ParseFiles(
		e.DirVhostHome+"/template"+"/"+tname+".html",
		e.DirVhostHome+"/template"+"/header.html",
		e.DirVhostHome+"/template"+"/sidebar.html",
		e.DirVhostHome+"/template"+"/footer.html",
	)
	if err != nil {
		e.printTmplPageError(err)
		return true
	}
	tmpl.Execute(*e.W, tmplDataAll{
		System: e.tmplGetSystemData(),
		Data:   data,
	})
	return true
}

func (e *Wrapper) TmplBackEnd(tcont []byte, data interface{}) bool {
	tmpl, err := template.New("template").Parse(string(tcont))
	if err != nil {
		e.printTmplPageError(err)
		return true
	}
	tmpl.Execute(*e.W, tmplDataAll{
		System: e.tmplGetSystemData(),
		Data:   data,
	})
	return true
}
