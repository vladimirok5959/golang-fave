package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type customHandler func(h http.Handler) http.Handler
type callbackBeforeAfter func(w http.ResponseWriter, r *http.Request, o interface{})
type callbackServer func(s *http.Server)

type bootstrap struct {
	path   string
	before callbackBeforeAfter
	after  callbackBeforeAfter
	object interface{}
}

func new(path string, before callbackBeforeAfter, after callbackBeforeAfter, object interface{}) *bootstrap {
	return &bootstrap{path, before, after, object}
}

func (this *bootstrap) handler(w http.ResponseWriter, r *http.Request) {
	if this.before != nil {
		this.before(w, r, this.object)
	}
	if r.URL.Path == "/"+this.path+"/bootstrap.css" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "text/css")
		w.Write(resource_bootstrap_css)
		return
	} else if r.URL.Path == "/"+this.path+"/bootstrap.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_bootstrap_js)
		return
	} else if r.URL.Path == "/"+this.path+"/jquery.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_jquery_js)
		return
	} else if r.URL.Path == "/"+this.path+"/popper.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_popper_js)
		return
	}
	if this.after != nil {
		this.after(w, r, this.object)
	}
}

func Start(h customHandler, host string, timeout time.Duration, path string, before callbackBeforeAfter, after callbackBeforeAfter, cbserv callbackServer, object interface{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", new(path, before, after, object).handler)

	var srv *http.Server
	if h == nil {
		srv = &http.Server{
			Addr:    host,
			Handler: mux,
		}
	} else {
		srv = &http.Server{
			Addr:    host,
			Handler: h(mux),
		}
	}

	if cbserv != nil {
		cbserv(srv)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go func() {
		fmt.Printf("Starting server at http://%s/\n", host)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				fmt.Println(err)
				stop <- os.Interrupt
				os.Exit(1)
			}
		}
	}()

	switch val := <-stop; val {
	case syscall.SIGTERM:
		fmt.Println("Shutting down server (terminate)...")
	case syscall.SIGINT:
		fmt.Println("Shutting down server (interrupt)...")
	default:
		fmt.Println("Shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
