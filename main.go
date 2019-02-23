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

var ParamHost string
var ParamPort int
var ParamWwwDir string

func init() {
	flag.StringVar(&ParamHost, "host", "0.0.0.0", "server host")
	flag.IntVar(&ParamPort, "port", 8080, "server port")
	flag.StringVar(&ParamWwwDir, "dir", "", "virtual hosts directory")
	flag.Parse()
}

func main() {
	// Init logger
	lg := logger.New()
	defer lg.Close()

	// Check www dir
	ParamWwwDir = utils.FixPath(ParamWwwDir)
	if !utils.IsDirExists(ParamWwwDir) {
		lg.Log("Virtual hosts directory is not exists", nil, true)
		lg.Log("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts", nil, true)
		return
	}

	// Attach www dir to logger
	lg.SetWwwDir(ParamWwwDir)

	// Session cleaner
	sess_cl_ch, sess_cl_stop := session_clean_start(ParamWwwDir)
	defer session_clean_stop(sess_cl_ch, sess_cl_stop)

	// Init mounted resources
	res := resource.New()
	res.Add(consts.AssetsCpScriptsJs, "application/javascript; charset=utf-8", assets.CpScriptsJs)
	res.Add(consts.AssetsCpStylesCss, "text/css", assets.CpStylesCss)
	res.Add(consts.AssetsSysBgPng, "image/png", assets.SysBgPng)
	res.Add(consts.AssetsSysFaveIco, "image/x-icon", assets.SysFaveIco)
	res.Add(consts.AssetsSysLogoPng, "image/png", assets.SysLogoPng)
	res.Add(consts.AssetsSysLogoSvg, "image/svg+xml", assets.SysLogoSvg)
	res.Add(consts.AssetsSysStylesCss, "text/css", assets.SysStylesCss)

	// Init static files helper
	stat := static.New(consts.DirIndexFile)

	// Init modules
	mods := modules.New()

	// Init and start web server
	bootstrap.Start(lg.Handler, fmt.Sprintf("%s:%d", ParamHost, ParamPort), 30, consts.AssetsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "fave.pro/"+consts.ServerVersion)
	}, func(w http.ResponseWriter, r *http.Request) {
		// Schema
		r.URL.Scheme = "http"

		// Mounted assets
		if res.Response(w, r, func(w http.ResponseWriter, r *http.Request, i *resource.Resource) {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}, nil) {
			return
		}

		// Host and port
		host, port := utils.ExtractHostPort(r.Host, false)
		curr_host := host
		vhost_dir := ParamWwwDir + string(os.PathSeparator) + host
		if !utils.IsDirExists(vhost_dir) {
			curr_host = "localhost"
			vhost_dir = ParamWwwDir + string(os.PathSeparator) + "localhost"
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
