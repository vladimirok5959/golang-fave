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

	"golang-fave/assets"
	"golang-fave/cblocks"
	"golang-fave/consts"
	"golang-fave/domains"
	"golang-fave/engine"
	"golang-fave/engine/basket"
	"golang-fave/engine/mysqlpool"
	"golang-fave/logger"
	"golang-fave/modules"
	"golang-fave/support"
	"golang-fave/utils"

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
	if err := support.New().Migration(consts.ParamWwwDir); err != nil {
		fmt.Printf("[ERROR] MIGRATION FAILED: %s\n", err)
	}

	// Init logger
	lg := logger.New()

	// Attach www dir to logger
	lg.SetWwwDir(consts.ParamWwwDir)

	// MySQL connections pool
	mpool := mysqlpool.New()

	// Session cleaner
	wSessCl := session_cleaner(consts.ParamWwwDir)

	// Image processing
	wImageGen := image_generator(consts.ParamWwwDir)

	// Xml generation
	wXmlGen := xml_generator(consts.ParamWwwDir, mpool)

	// Init mounted resources
	res := resource.New()
	assets.PopulateResources(res)

	// Init static files helper
	stat := static.New(consts.DirIndexFile)

	// Init modules
	mods := modules.New()

	// SMTP sender
	smtp_cl_ch, smtp_cl_stop := smtp_start(consts.ParamWwwDir, mpool)
	defer smtp_stop(smtp_cl_ch, smtp_cl_stop)

	// Shop basket
	sb := basket.New()
	sb_cl_ch, sb_cl_stop := basket_clean_start(sb)
	defer basket_clean_stop(sb_cl_ch, sb_cl_stop)

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
		var lg *logger.Logger
		if v, ok := (*o)[0].(*logger.Logger); ok {
			lg = v
		}

		var mpool *mysqlpool.MySqlPool
		if v, ok := (*o)[1].(*mysqlpool.MySqlPool); ok {
			mpool = v
		}

		var res *resource.Resource
		if v, ok := (*o)[2].(*resource.Resource); ok {
			res = v
		}

		var stat *static.Static
		if v, ok := (*o)[3].(*static.Static); ok {
			stat = v
		}

		var mods *modules.Modules
		if v, ok := (*o)[4].(*modules.Modules); ok {
			mods = v
		}
		// ---

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

		// Robots.txt and styles.css from template dir
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
				sb,
				lg,
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

		// ---
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

		if lg, ok := (*o)[0].(*logger.Logger); ok {
			lg.Close()
		}
		// ---

		if len(errs) > 0 {
			return errors.New("Shutdown callback: " + strings.Join(errs, ", "))
		}

		return nil
	}

	// Start server
	bootstrap.Start(
		&bootstrap.Opts{
			Handle:   lg.Handler,
			Host:     server_address,
			Path:     consts.AssetsPath,
			Cbserv:   server_params,
			Before:   before,
			After:    after,
			Timeout:  8 * time.Second,
			Shutdown: shutdown,
			Objects: &[]bootstrap.Iface{
				lg,
				mpool,
				wSessCl,
				wImageGen,
				wXmlGen,
				res,
				stat,
				mods,
			},
		},
	)
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
