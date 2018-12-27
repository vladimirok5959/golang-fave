package wrapper

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang-fave/engine/sessions"
)

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

type handleRun func(e *Wrapper) bool

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
	/*
		if e.redirectSeoFix() {
			e.Log("301")
			return
		}
	*/

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

	// Logic
	ret := false
	if hRun != nil {
		if hRun(e) {
			ret = true
		}
	}

	// Save session
	e.Session.Save()

	if ret {
		return
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
