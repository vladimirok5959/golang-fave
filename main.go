package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang-fave/engine"
	"golang-fave/engine/assets"
	"golang-fave/engine/basket"
	"golang-fave/engine/cblocks"
	"golang-fave/engine/consts"
	"golang-fave/engine/domains"
	"golang-fave/engine/logger"
	"golang-fave/engine/modules"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/utils"
	"golang-fave/engine/workers"
	"golang-fave/support"

	"github.com/vladimirok5959/golang-server-bootstrap/bootstrap"
	"github.com/vladimirok5959/golang-server-resources/resource"
	"github.com/vladimirok5959/golang-server-sessions/session"
	"github.com/vladimirok5959/golang-server-static/static"
	"github.com/vladimirok5959/golang-worker/worker"
)

func init() {
	flag.StringVar(&consts.ParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&consts.ParamPort, "port", 8080, "server port")
	flag.StringVar(&consts.ParamWwwDir, "dir", "", "virtual hosts directory")
	flag.BoolVar(&consts.ParamDebug, "debug", false, "debug mode with ignoring log files")
	flag.BoolVar(&consts.ParamKeepAlive, "keepalive", false, "enable/disable server keep alive")
	flag.Parse()
}

func main() {
	// Params from env vars
	read_env_params()

	// Check www dir
	consts.ParamWwwDir = utils.FixPath(consts.ParamWwwDir)
	if !utils.IsDirExists(consts.ParamWwwDir) {
		fmt.Printf("Virtual hosts directory is not exists\n")
		fmt.Printf("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts\n")
		return
	}

	// Run database migration
	if err := support.New().
		Migration(context.Background(), consts.ParamWwwDir); err != nil {
		fmt.Printf("[MIGRATION] FAILED: %s\n", err)
	} else {
		fmt.Printf("[MIGRATION] DONE!\n")
	}

	// Init logger
	logs := logger.New()

	// Attach www dir to logger
	logs.SetWwwDir(consts.ParamWwwDir)

	// MySQL connections pool
	mpool := mysqlpool.New()

	// Session cleaner
	wSessCl := workers.SessionCleaner(consts.ParamWwwDir)

	// Image processing
	wImageGen := workers.ImageGenerator(consts.ParamWwwDir)

	// Xml generation
	wXmlGen := workers.XmlGenerator(consts.ParamWwwDir, mpool)

	// SMTP sender
	wSmtpSnd := workers.SmtpSender(consts.ParamWwwDir, mpool)

	// Init mounted resources
	res := resource.New()
	assets.PopulateResources(res)

	// Init static files helper
	stat := static.New(consts.DirIndexFile)

	// Init modules
	mods := modules.New()

	// Shop basket
	shopBasket := basket.New()
	wBasketCl := workers.BasketCleaner(shopBasket)

	// Init cache blocks
	cbs := cblocks.New()

	// Init and start web server
	server_address := fmt.Sprintf("%s:%d", consts.ParamHost, consts.ParamPort)

	// Server params
	server_params := func(s *http.Server) {
		s.SetKeepAlivesEnabled(consts.ParamKeepAlive)
	}

	// Before callback
	before := func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		o *[]bootstrap.Iface,
	) {
		w.Header().Set("Server", "fave.pro/"+consts.ServerVersion)
	}

	// After callback
	after := func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		o *[]bootstrap.Iface,
	) {
		// Schema
		r.URL.Scheme = "http"

		// Convert
		var logs *logger.Logger
		if v, ok := (*o)[0].(*logger.Logger); ok {
			logs = v
		}

		var mpool *mysqlpool.MySqlPool
		if v, ok := (*o)[1].(*mysqlpool.MySqlPool); ok {
			mpool = v
		}

		var res *resource.Resource
		if v, ok := (*o)[6].(*resource.Resource); ok {
			res = v
		}

		var stat *static.Static
		if v, ok := (*o)[7].(*static.Static); ok {
			stat = v
		}

		var mods *modules.Modules
		if v, ok := (*o)[8].(*modules.Modules); ok {
			mods = v
		}

		var shopBasket *basket.Basket
		if v, ok := (*o)[9].(*basket.Basket); ok {
			shopBasket = v
		}

		// Mounted assets
		if res.Response(
			w,
			r,
			func(
				w http.ResponseWriter,
				r *http.Request,
				i *resource.OneResource,
			) {
				if consts.ParamDebug && i.Path == "assets/cp/scripts.js" {
					w.Write([]byte("window.fave_debug=true;"))
				}
			},
			nil,
		) {
			return
		}

		// Host and port
		host, port := utils.ExtractHostPort(r.Host, false)
		curr_host := host

		// Domain bindings
		doms := domains.New(consts.ParamWwwDir)
		if mhost := doms.GetHost(host); mhost != "" {
			curr_host = mhost
		}

		vhost_dir := consts.ParamWwwDir + string(os.PathSeparator) + curr_host
		if !utils.IsDirExists(vhost_dir) {
			curr_host = "localhost"
			vhost_dir = consts.ParamWwwDir + string(os.PathSeparator) + "localhost"
		}

		// Check for minimal dirs structure
		vhost_dir_config := vhost_dir + string(os.PathSeparator) + "config"
		vhost_dir_htdocs := vhost_dir + string(os.PathSeparator) + "htdocs"
		vhost_dir_logs := vhost_dir + string(os.PathSeparator) + "logs"
		vhost_dir_template := vhost_dir + string(os.PathSeparator) + "template"
		vhost_dir_tmp := vhost_dir + string(os.PathSeparator) + "tmp"
		if !utils.IsDirExists(vhost_dir_config) {
			utils.SystemErrorPageEngine(
				w,
				errors.New("Folder "+vhost_dir_config+" is not found"),
			)
			return
		}
		if !utils.IsDirExists(vhost_dir_htdocs) {
			utils.SystemErrorPageEngine(
				w,
				errors.New("Folder "+vhost_dir_htdocs+" is not found"),
			)
			return
		}
		if !utils.IsDirExists(vhost_dir_logs) {
			utils.SystemErrorPageEngine(
				w,
				errors.New("Folder "+vhost_dir_logs+" is not found"),
			)
			return
		}
		if !utils.IsDirExists(vhost_dir_template) {
			utils.SystemErrorPageEngine(
				w,
				errors.New("Folder "+vhost_dir_template+" is not found"),
			)
			return
		}
		if !utils.IsDirExists(vhost_dir_tmp) {
			utils.SystemErrorPageEngine(
				w,
				errors.New("Folder "+vhost_dir_tmp+" is not found"),
			)
			return
		}

		// Static files
		if stat.Response(vhost_dir_htdocs, w, r, nil, nil) {
			return
		}

		// robots.txt, styles.css and scripts.js from templates dir
		if ServeTemplateFile(w, r, "robots.txt", "", vhost_dir_template) {
			return
		}
		if ServeTemplateFile(w, r, "styles.css", "assets/theme/", vhost_dir_template) {
			return
		}
		if ServeTemplateFile(w, r, "scripts.js", "assets/theme/", vhost_dir_template) {
			return
		}

		// Session
		sess := session.New(w, r, vhost_dir_tmp)
		defer sess.Close()

		// Logic
		if mpool != nil {
			if engine.Response(
				mpool,
				shopBasket,
				logs,
				mods,
				w,
				r,
				sess,
				cbs,
				host,
				port,
				curr_host,
				vhost_dir_config,
				vhost_dir_htdocs,
				vhost_dir_logs,
				vhost_dir_template,
				vhost_dir_tmp,
			) {
				return
			}
		}

		// Error 404
		utils.SystemErrorPage404(w)
	}

	// Shutdown callback
	shutdown := func(
		ctx context.Context,
		o *[]bootstrap.Iface,
	) error {
		var errs []string

		if wBasketCl, ok := (*o)[10].(*worker.Worker); ok {
			if err := wBasketCl.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", wBasketCl, err.Error()))
			}
		}

		if wSmtpSnd, ok := (*o)[5].(*worker.Worker); ok {
			if err := wSmtpSnd.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", wSmtpSnd, err.Error()))
			}
		}

		if wXmlGen, ok := (*o)[4].(*worker.Worker); ok {
			if err := wXmlGen.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", wXmlGen, err.Error()))
			}
		}

		if wImageGen, ok := (*o)[3].(*worker.Worker); ok {
			if err := wImageGen.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", wImageGen, err.Error()))
			}
		}

		if wSessCl, ok := (*o)[2].(*worker.Worker); ok {
			if err := wSessCl.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", wSessCl, err.Error()))
			}
		}

		if mpool, ok := (*o)[1].(*mysqlpool.MySqlPool); ok {
			if err := mpool.Close(); err != nil {
				errs = append(errs, fmt.Sprintf("(%T): %s", mpool, err.Error()))
			}
		}

		if logs, ok := (*o)[0].(*logger.Logger); ok {
			logs.Close()
		}

		if len(errs) > 0 {
			return errors.New("Shutdown callback: " + strings.Join(errs, ", "))
		}

		return nil
	}

	// Start server
	bootstrap.Start(
		&bootstrap.Opts{
			Handle:   logs.Handler,
			Host:     server_address,
			Path:     consts.AssetsPath,
			Cbserv:   server_params,
			Before:   before,
			After:    after,
			Timeout:  8 * time.Second,
			Shutdown: shutdown,
			Objects: &[]bootstrap.Iface{
				logs,
				mpool,
				wSessCl,
				wImageGen,
				wXmlGen,
				wSmtpSnd,
				res,
				stat,
				mods,
				shopBasket,
				wBasketCl,
			},
		},
	)
}

func read_env_params() {
	if consts.ParamHost == "0.0.0.0" {
		if os.Getenv("FAVE_HOST") != "" {
			consts.ParamHost = os.Getenv("FAVE_HOST")
		}
	}
	if consts.ParamPort == 8080 {
		if os.Getenv("FAVE_PORT") != "" {
			consts.ParamPort = utils.StrToInt(os.Getenv("FAVE_PORT"))
		}
	}
	if consts.ParamWwwDir == "" {
		if os.Getenv("FAVE_DIR") != "" {
			consts.ParamWwwDir = os.Getenv("FAVE_DIR")
		}
	}
	if consts.ParamDebug == false {
		if os.Getenv("FAVE_DEBUG") == "true" {
			consts.ParamDebug = true
		}
	}
	if consts.ParamKeepAlive == false {
		if os.Getenv("FAVE_KEEPALIVE") == "true" {
			consts.ParamKeepAlive = true
		}
	}
}

func ServeTemplateFile(
	w http.ResponseWriter,
	r *http.Request,
	file string,
	path string,
	dir string,
) bool {
	if r.URL.Path == "/"+path+file {
		if utils.IsRegularFileExists(dir + string(os.PathSeparator) + file) {
			http.ServeFile(w, r, dir+string(os.PathSeparator)+file)
			return true
		}
	}
	return false
}
