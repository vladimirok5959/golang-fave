package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang-fave/consts"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-bootstrap/bootstrap"
	/*
		"context"
		"errors"
		"flag"
		"fmt"
		"log"
		"net/http"
		"os"
		"os/signal"
		"strconv"
		"strings"
		"time"

		"golang-fave/constants"
		"golang-fave/engine/actions"
		"golang-fave/engine/wrapper"
	*/)

var ParamHost string
var ParamPort int
var ParamWwwDir string

//var VhostHomeDir string

func init() {
	flag.StringVar(&ParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&ParamPort, "port", 8080, "server port")
	flag.StringVar(&ParamWwwDir, "dir", "", "virtual hosts directory")
	flag.Parse()
}

func main() {
	ParamWwwDir = utils.FixPath(ParamWwwDir)
	if !(utils.IsFileExists(ParamWwwDir) && utils.IsDir(ParamWwwDir)) {
		fmt.Println("Virtual hosts directory is not exists")
		fmt.Println("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts")
		return
	}

	bootstrap.Start("127.0.0.1:8080", 30, "assets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "fave.pro/"+consts.ServerVersion)
	}, func(w http.ResponseWriter, r *http.Request) {

		/*
			// After callback
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<div>Hello World!</div>
				<div><a href="/assets/bootstrap.css">/assets/bootstrap.css</a></div>
				<div><a href="/assets/bootstrap.js">/assets/bootstrap.js</a></div>
				<div><a href="/assets/jquery.js">/assets/jquery.js</a></div>
				<div><a href="/assets/popper.js">/assets/popper.js</a></div>
			`))
		*/

	})

	/*
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
	*/
}

/*
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
	w.Header().Set("Server", "fave.pro/"+constants.ServerVersion)

	wrap := wrapper.New(&w, r, host, port, ParamWwwDir, VhostHomeDir)

	// Check if localhost exists
	if !vhExists(VhostHomeDir) && !strings.HasPrefix(r.URL.Path, "/assets/") {
		wrap.PrintEnginePageError(errors.New("Folder " + VhostHomeDir + " is not found"))
		return
	}

	// Create and start engine
	wrap.Run(func(wrapper *wrapper.Wrapper) bool {
		// Actions
		action := actions.New(wrapper)
		if action.Run() {
			wrapper.Session.Save()
			return true
		}

		// Pages
		return handlerPage(wrapper)
	})
}
*/
