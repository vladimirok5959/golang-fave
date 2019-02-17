package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type hndl func(h http.Handler) http.Handler
type callback func(w http.ResponseWriter, r *http.Request)

type bootstrap struct {
	path   string
	before callback
	after  callback
}

func new(path string, before callback, after callback) *bootstrap {
	return &bootstrap{path, before, after}
}

func (this *bootstrap) handler(w http.ResponseWriter, r *http.Request) {
	if this.before != nil {
		this.before(w, r)
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
		this.after(w, r)
	}
}

func Start(h hndl, host string, timeout time.Duration, path string, before callback, after callback) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", new(path, before, after).handler)

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

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		fmt.Printf("Starting server at http://%s/\n", host)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				fmt.Println(err)
			}
		}
	}()
	<-stop
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
