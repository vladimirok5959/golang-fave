package engine

import (
	"net/http"
	"os"
	"strings"

	"golang-fave/engine/assets"
	"golang-fave/engine/basket"
	"golang-fave/engine/cblocks"
	"golang-fave/engine/logger"
	"golang-fave/engine/modules"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"

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
	// Request was canceled
	if this.contextDone() {
		return false
	}

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

	// Request was canceled
	if this.contextDone() {
		return false
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

	// Request was canceled
	if this.contextDone() {
		return false
	}

	// Check for MySQL connection
	err := this.Wrap.UseDatabase()
	if err != nil {
		utils.SystemErrorPageEngine(this.Wrap.W, err)
		return true
	}

	// Request was canceled
	if this.contextDone() {
		return false
	}

	// Separated logic
	if !this.Wrap.IsBackend {
		// Maintenance mode
		if this.Wrap.Config.Engine.Maintenance != 0 {
			if this.Wrap.User == nil {
				// this.Wrap.UseDatabase()
				this.Wrap.LoadSessionUser()
			}
			if this.Wrap.User == nil {
				this.Wrap.RenderFrontEnd("maintenance", nil, http.StatusServiceUnavailable)
				return true
			}
			if this.Wrap.User.A_id <= 0 {
				this.Wrap.RenderFrontEnd("maintenance", nil, http.StatusServiceUnavailable)
				return true
			}
			if this.Wrap.User.A_admin <= 0 {
				this.Wrap.RenderFrontEnd("maintenance", nil, http.StatusServiceUnavailable)
				return true
			}
		}

		// Render frontend
		return this.Mods.XXXFrontEnd(this.Wrap)
	}

	// Request was canceled
	if this.contextDone() {
		return false
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

	// Request was canceled
	if this.contextDone() {
		return false
	}

	// Try load current user data
	if !this.Wrap.LoadSessionUser() {
		http.Redirect(this.Wrap.W, this.Wrap.R, "/", 302)
		return true
	}

	// Request was canceled
	if this.contextDone() {
		return false
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

	// Request was canceled
	if this.contextDone() {
		return false
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

func (this *Engine) contextDone() bool {
	select {
	case <-this.Wrap.R.Context().Done():
		return true
	default:
	}
	return false
}
