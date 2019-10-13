package engine

import (
	"net/http"
	"os"
	"strings"

	"golang-fave/assets"
	"golang-fave/cblocks"
	"golang-fave/engine/basket"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/wrapper"
	"golang-fave/logger"
	"golang-fave/modules"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Engine struct {
	Wrap *wrapper.Wrapper
	Mods *modules.Modules
}

func Response(mp *mysqlpool.MySqlPool, sb *basket.Basket, l *logger.Logger, m *modules.Modules, w http.ResponseWriter, r *http.Request, s *session.Session, c *cblocks.CacheBlocks, host, port, chost, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) bool {
	wrap := wrapper.New(l, w, r, s, c, host, port, chost, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp, mp, sb)
	eng := &Engine{
		Wrap: wrap,
		Mods: m,
	}
	return eng.Process()
}

func (this *Engine) Process() bool {
	this.Wrap.IsBackend = this.Wrap.R.URL.Path == "/cp" || strings.HasPrefix(this.Wrap.R.URL.Path, "/cp/")
	this.Wrap.ConfMysqlExists = utils.IsMySqlConfigExists(this.Wrap.DConfig + string(os.PathSeparator) + "mysql.json")
	this.Wrap.UrlArgs = append(this.Wrap.UrlArgs, utils.UrlToArray(this.Wrap.R.URL.Path)...)
	if this.Wrap.IsBackend && len(this.Wrap.UrlArgs) >= 1 && this.Wrap.UrlArgs[0] == "cp" {
		this.Wrap.UrlArgs = this.Wrap.UrlArgs[1:]
	}

	// Action
	if this.Mods.XXXActionFire(this.Wrap) {
		return true
	}

	// Redirect to CP for creating MySQL config file
	if !this.Wrap.IsBackend && !this.Wrap.ConfMysqlExists {
		this.Wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(this.Wrap.W, this.Wrap.R, this.Wrap.R.URL.Scheme+"://"+this.Wrap.R.Host+"/cp/", 302)
		return true
	}

	// Display MySQL install page on backend
	if this.Wrap.IsBackend && !this.Wrap.ConfMysqlExists {
		// Redirect
		if this.redirectFixCpUrl() {
			return true
		}
		// Show mysql settings form
		utils.SystemRenderTemplate(this.Wrap.W, assets.TmplCpMySql, nil, "", "")
		return true
	}

	// Separated logic
	if !this.Wrap.IsBackend {
		// Render frontend
		return this.Mods.XXXFrontEnd(this.Wrap)
	}

	// Backend must use MySQL anyway, so, check and connect
	err := this.Wrap.UseDatabase()
	if err != nil {
		utils.SystemErrorPageEngine(this.Wrap.W, err)
		return true
	}

	// Show login page if need
	if this.Wrap.S.GetInt("UserId", 0) <= 0 {
		// Redirect
		if this.redirectFixCpUrl() {
			return true
		}
		// Show login form
		utils.SystemRenderTemplate(this.Wrap.W, assets.TmplCpLogin, nil, "", "")
		return true
	}

	// Try load current user data
	if !this.Wrap.LoadSessionUser() {
		http.Redirect(this.Wrap.W, this.Wrap.R, "/", 302)
		return true
	}

	// Only active admins can use backend
	if !(this.Wrap.User.A_admin == 1 && this.Wrap.User.A_active == 1) {
		// Redirect
		if this.redirectFixCpUrl() {
			return true
		}
		// Show login form
		utils.SystemRenderTemplate(this.Wrap.W, assets.TmplCpLogin, nil, "", "")
		return true
	}

	// Redirect
	if this.redirectFixCpUrl() {
		return true
	}

	// Render backend
	return this.Mods.XXXBackEnd(this.Wrap)
}

func (this *Engine) redirectFixCpUrl() bool {
	if len(this.Wrap.R.URL.Path) > 0 && this.Wrap.R.URL.Path[len(this.Wrap.R.URL.Path)-1] != '/' {
		http.Redirect(this.Wrap.W, this.Wrap.R, this.Wrap.R.URL.Path+"/"+utils.ExtractGetParams(this.Wrap.R.RequestURI), 302)
		return true
	}
	return false
}
