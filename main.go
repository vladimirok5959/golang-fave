package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang-fave/assets"
	"golang-fave/consts"
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
	ParamWwwDir = utils.FixPath(ParamWwwDir)
	if !utils.IsHostDirExists(ParamWwwDir) {
		fmt.Println("Virtual hosts directory is not exists")
		fmt.Println("Example: ./fave -host 127.0.0.1 -port 80 -dir ./hosts")
		return
	}

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

	// TODO: Logic as object here
	// Init logic

	bootstrap.Start(fmt.Sprintf("%s:%d", ParamHost, ParamPort), 30, consts.AssetsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "fave.pro/"+consts.ServerVersion)
	}, func(w http.ResponseWriter, r *http.Request) {
		// Mounted assets
		if res.Response(w, r, func(w http.ResponseWriter, r *http.Request, i *resource.Resource) {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}, nil) {
			return
		}

		// Host and port
		host, port := utils.ExtractHostPort(r.Host, false)
		vhost_dir := ParamWwwDir + string(os.PathSeparator) + host
		if !utils.IsHostDirExists(vhost_dir) {
			vhost_dir = ParamWwwDir + string(os.PathSeparator) + "localhost"
		}

		// Check for minimal dir structure
		vhost_dir_config := vhost_dir + string(os.PathSeparator) + "config"
		vhost_dir_htdocs := vhost_dir + string(os.PathSeparator) + "htdocs"
		vhost_dir_logs := vhost_dir + string(os.PathSeparator) + "logs"
		vhost_dir_template := vhost_dir + string(os.PathSeparator) + "template"
		vhost_dir_tmp := vhost_dir + string(os.PathSeparator) + "tmp"

		if !utils.IsHostDirExists(vhost_dir_config) {
			utils.SystemErrorPage(w, errors.New("Folder "+vhost_dir_config+" is not found"))
			return
		}
		if !utils.IsHostDirExists(vhost_dir_htdocs) {
			utils.SystemErrorPage(w, errors.New("Folder "+vhost_dir_htdocs+" is not found"))
			return
		}
		if !utils.IsHostDirExists(vhost_dir_logs) {
			utils.SystemErrorPage(w, errors.New("Folder "+vhost_dir_logs+" is not found"))
			return
		}
		if !utils.IsHostDirExists(vhost_dir_template) {
			utils.SystemErrorPage(w, errors.New("Folder "+vhost_dir_template+" is not found"))
			return
		}
		if !utils.IsHostDirExists(vhost_dir_tmp) {
			utils.SystemErrorPage(w, errors.New("Folder "+vhost_dir_tmp+" is not found"))
			return
		}

		// Static files
		if stat.Response(vhost_dir_htdocs, w, r, nil, nil) {
			return
		}

		// Session
		//sess := session.New(w, r, vhost_dir_tmp)
		//defer sess.Close()

		// Session struct need to make public!
		sess := session.New(w, r, vhost_dir_tmp)
		defer sess.Close()

		// Logic
		// TODO: call from `logic.Response()`
		// TODO: create logic object here???
		if logic(host, port, vhost_dir_config, vhost_dir_htdocs, vhost_dir_logs, vhost_dir_template, vhost_dir_tmp, w, r) {
			return
		}

		// Error 404
		// TODO: display default template
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`404`))
	})

	// Delete expired session files
	// session.Clean("./tmp")
}

func logic(host, port, dir_config, dir_htdocs, dir_logs, dir_template, dir_tmp string, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "text/html")

		//counter := sess.GetInt("counter", 0)
		//w.Write([]byte(`Logic -> (` + fmt.Sprintf("%d", counter) + `)`))

		w.Write([]byte(`Logic`))

		//counter++
		//sess.SetInt("counter", counter)

		return true
	}
	return false
}
