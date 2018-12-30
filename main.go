package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/actions"
	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

const C_Debug = !false

var ParamHost string
var ParamPort int
var ParamWwwDir string
var VhostHomeDir string

func init() {
	flag.StringVar(&ParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&ParamPort, "port", 8080, "server port")
	flag.StringVar(&ParamWwwDir, "dir", "", "virtual hosts directory")
	flag.Parse()
}

func main() {
	if _, err := os.Stat(ParamWwwDir); os.IsNotExist(err) {
		fmt.Println("Virtual hosts directory is not exists")
		fmt.Println("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts")
		return
	}
	if ParamWwwDir[len(ParamWwwDir)-1] != '/' {
		ParamWwwDir = ParamWwwDir + "/"
	}

	// Handle
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ParamHost, ParamPort),
		Handler: mux,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server at %s:%d", ParamHost, ParamPort)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// Wait for signal
	<-stop

	log.Printf("Shutting down server...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server is off!")
	}
}

func vhExists(vhosthome string) bool {
	if st, err := os.Stat(vhosthome); !os.IsNotExist(err) {
		if err == nil {
			fmode := st.Mode()
			if fmode.IsDir() {
				return true
			}
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Build vhost home dir
	host := r.Host
	port := strconv.Itoa(ParamPort)
	index := strings.Index(host, ":")
	if index > -1 {
		port = host[index+1:]
		host = host[0:index]
	}

	// Cut "www" if exists
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}
	VhostHomeDir = ParamWwwDir + host

	// Check if virtual host exists
	if !vhExists(VhostHomeDir) {
		host = "localhost"
		VhostHomeDir = ParamWwwDir + host
	}

	// Set protocol
	r.URL.Scheme = "http"

	// Set server name
	w.Header().Set("Server", "fave.pro")

	// Create and start engine
	wrapper.New(&w, r, host, port, ParamWwwDir, VhostHomeDir, C_Debug).
		Run(func(wrapper *wrapper.Wrapper) bool {
			// Actions
			action := actions.New(wrapper)
			if action.Call() {
				wrapper.Session.Save()
				return true
			}

			// Pages
			return handlerPage(wrapper)
		})
}

func handlerPage(wrapper *wrapper.Wrapper) bool {
	if !(wrapper.R.URL.Path == "/cp" || strings.HasPrefix(wrapper.R.URL.Path, "/cp/")) {
		return handlerFrontEnd(wrapper)
	} else {
		return handlerBackEnd(wrapper)
	}
}

func handlerFrontEnd(wrapper *wrapper.Wrapper) bool {
	// Redirect to CP, if MySQL config file is not exists
	if !utils.IsMySqlConfigExists(wrapper.DirVHostHome) {
		(*wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(*wrapper.W, wrapper.R, wrapper.R.URL.Scheme+"://"+wrapper.R.Host+"/cp/", 302)
		return true
	}

	// Connect to database

	// Else logic here
	if wrapper.R.URL.Path == "/" {
		return wrapper.TmplFrontEnd("index", nil)
	}

	return false
}

func handlerBackEnd(wrapper *wrapper.Wrapper) bool {
	// MySQL config page
	if !utils.IsMySqlConfigExists(wrapper.DirVHostHome) {
		return wrapper.TmplBackEnd(templates.CpMySQL, nil)
	}

	// Connect to database
	mc, err := utils.MySqlConfigRead(wrapper.DirVHostHome)
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}
	db, err := sql.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}
	defer db.Close()
	err = db.Ping()
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}

	// Check if any user exists

	// Login page
	return wrapper.TmplBackEnd(templates.CpLogin, nil)
}
