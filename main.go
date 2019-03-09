package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine"
	"golang-fave/logger"
	"golang-fave/modules"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-bootstrap/bootstrap"
	"github.com/vladimirok5959/golang-server-resources/resource"
	"github.com/vladimirok5959/golang-server-sessions/session"
	"github.com/vladimirok5959/golang-server-static/static"
)

func init() {
	flag.StringVar(&consts.ParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&consts.ParamPort, "port", 8080, "server port")
	flag.StringVar(&consts.ParamWwwDir, "dir", "", "virtual hosts directory")
	flag.BoolVar(&consts.ParamDebug, "debug", false, "debug mode with ignoring log files")
	flag.Parse()
}

func main() {
	// Universal, params by env vars
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

	// Check www dir
	consts.ParamWwwDir = utils.FixPath(consts.ParamWwwDir)
	if !utils.IsDirExists(consts.ParamWwwDir) {
		fmt.Printf("Virtual hosts directory is not exists\n")
		fmt.Printf("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts\n")
		return
	}

	// Init logger
	lg := logger.New()
	defer lg.Close()

	// Attach www dir to logger
	lg.SetWwwDir(consts.ParamWwwDir)

	// Session cleaner
	sess_cl_ch, sess_cl_stop := session_clean_start(consts.ParamWwwDir)
	defer session_clean_stop(sess_cl_ch, sess_cl_stop)

	// Init mounted resources
	res := resource.New()
	assets.PopulateResources(res)

	// Init static files helper
	stat := static.New(consts.DirIndexFile)

	// Init modules
	mods := modules.New()

	// Init and start web server
	bootstrap.Start(lg.Handler, fmt.Sprintf("%s:%d", consts.ParamHost, consts.ParamPort), 9, consts.AssetsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "fave.pro/"+consts.ServerVersion)
	}, func(w http.ResponseWriter, r *http.Request) {
		// Schema
		r.URL.Scheme = "http"

		// Mounted assets
		if res.Response(w, r, func(w http.ResponseWriter, r *http.Request, i *resource.OneResource) {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}, nil) {
			return
		}

		// Host and port
		host, port := utils.ExtractHostPort(r.Host, false)
		curr_host := host
		vhost_dir := consts.ParamWwwDir + string(os.PathSeparator) + host
		if !utils.IsDirExists(vhost_dir) {
			curr_host = "localhost"
			vhost_dir = consts.ParamWwwDir + string(os.PathSeparator) + "localhost"
		}

		// Check for minimal dir structure
		vhost_dir_config := vhost_dir + string(os.PathSeparator) + "config"
		vhost_dir_htdocs := vhost_dir + string(os.PathSeparator) + "htdocs"
		vhost_dir_logs := vhost_dir + string(os.PathSeparator) + "logs"
		vhost_dir_template := vhost_dir + string(os.PathSeparator) + "template"
		vhost_dir_tmp := vhost_dir + string(os.PathSeparator) + "tmp"
		if !utils.IsDirExists(vhost_dir_config) {
			utils.SystemErrorPageEngine(w, errors.New("Folder "+vhost_dir_config+" is not found"))
			return
		}
		if !utils.IsDirExists(vhost_dir_htdocs) {
			utils.SystemErrorPageEngine(w, errors.New("Folder "+vhost_dir_htdocs+" is not found"))
			return
		}
		if !utils.IsDirExists(vhost_dir_logs) {
			utils.SystemErrorPageEngine(w, errors.New("Folder "+vhost_dir_logs+" is not found"))
			return
		}
		if !utils.IsDirExists(vhost_dir_template) {
			utils.SystemErrorPageEngine(w, errors.New("Folder "+vhost_dir_template+" is not found"))
			return
		}
		if !utils.IsDirExists(vhost_dir_tmp) {
			utils.SystemErrorPageEngine(w, errors.New("Folder "+vhost_dir_tmp+" is not found"))
			return
		}

		// Static files
		if stat.Response(vhost_dir_htdocs, w, r, nil, nil) {
			return
		}

		// Session
		sess := session.New(w, r, vhost_dir_tmp)
		defer sess.Close()

		// Logic
		if engine.Response(lg, mods, w, r, sess, host, port, curr_host, vhost_dir_config, vhost_dir_htdocs, vhost_dir_logs, vhost_dir_template, vhost_dir_tmp) {
			return
		}

		// Error 404
		utils.SystemErrorPage404(w)
	})
}
