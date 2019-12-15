package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vladimirok5959/golang-ctrlc/ctrlc"
)

type customHandler func(h http.Handler) http.Handler
type callbackBeforeAfter func(ctx context.Context, w http.ResponseWriter, r *http.Request, o interface{})
type callbackServer func(s *http.Server)

type Opts struct {
	Handle  customHandler
	Host    string
	Timeout time.Duration
	Path    string
	Before  callbackBeforeAfter
	After   callbackBeforeAfter
	Cbserv  callbackServer
	Object  interface{}
}

type bootstrap struct {
	ctx  context.Context
	opts *Opts
}

func new(ctx context.Context, opts *Opts) *bootstrap {
	return &bootstrap{ctx: ctx, opts: opts}
}

func (this *bootstrap) handler(w http.ResponseWriter, r *http.Request) {
	if this.opts.Before != nil {
		this.opts.Before(this.ctx, w, r, this.opts.Object)
	}
	if r.URL.Path == "/"+this.opts.Path+"/bootstrap.css" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "text/css")
		w.Write(resource_bootstrap_css)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/bootstrap.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_bootstrap_js)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/jquery.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_jquery_js)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/popper.js" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(resource_popper_js)
		return
	}
	if this.opts.After != nil {
		this.opts.After(this.ctx, w, r, this.opts.Object)
	}
}

func Start(opts *Opts) {
	if opts == nil {
		fmt.Println("Start: options is not defined")
		os.Exit(1)
	}

	ctrlc.App(
		opts.Timeout,
		func(ctx context.Context, shutdown context.CancelFunc) *[]ctrlc.Iface {
			mux := http.NewServeMux()
			mux.HandleFunc("/", new(ctx, opts).handler)

			var srv *http.Server
			if opts.Handle == nil {
				srv = &http.Server{
					Addr:    opts.Host,
					Handler: mux,
				}
			} else {
				srv = &http.Server{
					Addr:    opts.Host,
					Handler: opts.Handle(mux),
				}
			}

			if opts.Cbserv != nil {
				opts.Cbserv(srv)
			}

			go func() {
				fmt.Printf("Starting server at http://%s/\n", opts.Host)
				if err := srv.ListenAndServe(); err != nil {
					if err != http.ErrServerClosed {
						fmt.Printf("Web server startup error: %s\n", err.Error())
						shutdown()
						return
					}
				}
			}()

			return &[]ctrlc.Iface{srv}
		},
	)
}
