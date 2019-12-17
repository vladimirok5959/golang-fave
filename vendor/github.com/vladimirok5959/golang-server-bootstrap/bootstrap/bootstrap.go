package bootstrap

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vladimirok5959/golang-ctrlc/ctrlc"
)

type Handler func(h http.Handler) http.Handler

type CBServer func(s *http.Server)

type Iface interface{}

type BeforeAfter func(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	o *[]Iface,
)

type ShutdownFunc func(ctx context.Context, o *[]Iface) error

type Opts struct {
	Handle   Handler
	Host     string
	Path     string
	Cbserv   CBServer
	Before   BeforeAfter
	After    BeforeAfter
	Objects  *[]Iface
	Timeout  time.Duration
	Shutdown ShutdownFunc
}

type bootstrap struct {
	ctx  context.Context
	opts *Opts
}

func new(ctx context.Context, opts *Opts) *bootstrap {
	return &bootstrap{ctx: ctx, opts: opts}
}

func etag(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func modified(p string, s int, v int64, w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Length", fmt.Sprintf("%d", s))
	w.Header().Set("Cache-Control", "no-cache")

	// Set: ETag
	ehash := etag(fmt.Sprintf("%s-%d-%d", p, s, v))
	w.Header().Set("ETag", fmt.Sprintf("%s", ehash))

	// Set: Last-Modified
	w.Header().Set(
		"Last-Modified",
		time.Unix(v, 0).In(time.FixedZone("GMT", 0)).Format("Wed, 01 Oct 2006 15:04:05 GMT"),
	)

	// Check: ETag
	if cc := r.Header.Get("Cache-Control"); cc != "no-cache" {
		if inm := r.Header.Get("If-None-Match"); inm == ehash {
			w.WriteHeader(http.StatusNotModified)
			return false
		}
	}

	// Check: Last-Modified
	if cc := r.Header.Get("Cache-Control"); cc != "no-cache" {
		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if t, err := time.Parse("Wed, 01 Oct 2006 15:04:05 GMT", ims); err == nil {
				if time.Unix(v, 0).In(time.FixedZone("GMT", 0)).Unix() <= t.In(time.FixedZone("GMT", 0)).Unix() {
					w.WriteHeader(http.StatusNotModified)
					return false
				}
			}
		}
	}

	return true
}

func (this *bootstrap) handler(w http.ResponseWriter, r *http.Request) {
	if this.opts.Before != nil {
		this.opts.Before(this.ctx, w, r, this.opts.Objects)
	}
	if r.URL.Path == "/"+this.opts.Path+"/bootstrap.css" {
		w.Header().Set("Content-Type", "text/css")
		if !modified(r.URL.Path, len(rbc), rbcm, w, r) {
			return
		}
		w.Write(rbc)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/bootstrap.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		if !modified(r.URL.Path, len(rbj), rbjm, w, r) {
			return
		}
		w.Write(rbj)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/jquery.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		if !modified(r.URL.Path, len(rjj), rjjm, w, r) {
			return
		}
		w.Write(rjj)
		return
	} else if r.URL.Path == "/"+this.opts.Path+"/popper.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		if !modified(r.URL.Path, len(rpj), rpjm, w, r) {
			return
		}
		w.Write(rpj)
		return
	}
	if this.opts.After != nil {
		this.opts.After(this.ctx, w, r, this.opts.Objects)
	}
}

func (this *bootstrap) Shutdown(ctx context.Context) error {
	if this.opts.Shutdown != nil {
		return this.opts.Shutdown(ctx, this.opts.Objects)
	}
	return nil
}

func Start(opts *Opts) {
	if opts == nil {
		fmt.Println("Start: options is not defined")
		os.Exit(1)
	}

	ctrlc.App(
		opts.Timeout,
		func(ctx context.Context, shutdown context.CancelFunc) *[]ctrlc.Iface {
			bt := new(ctx, opts)

			mux := http.NewServeMux()
			mux.HandleFunc("/", bt.handler)

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

			return &[]ctrlc.Iface{bt, srv}
		},
	)
}
