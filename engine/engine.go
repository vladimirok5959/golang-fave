package engine

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang-fave/assets"
	"golang-fave/engine/wrapper"
	"golang-fave/logger"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Engine struct {
	Wrap *wrapper.Wrapper
	// Actions
	// Front-end or Back-end
}

func Response(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) bool {
	wrap := wrapper.New(l, w, r, s, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp)
	eng := &Engine{Wrap: wrap}

	return eng.Process()
}

func (this *Engine) Process() bool {
	this.Wrap.IsBackend = this.Wrap.R.URL.Path == "/cp" || strings.HasPrefix(this.Wrap.R.URL.Path, "/cp/")
	this.Wrap.ConfMysqlExists = utils.IsMySqlConfigExists(this.Wrap.DConfig + string(os.PathSeparator) + "mysql.json")

	// Some system actions here...
	// MySQL install
	// MySQL first user
	// User login
	// User logout

	// Redirect to CP for creating MySQL config file
	if !this.Wrap.IsBackend && !this.Wrap.ConfMysqlExists {
		this.Wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(this.Wrap.W, this.Wrap.R, this.Wrap.R.URL.Scheme+"://"+this.Wrap.R.Host+"/cp/", 302)
		return true
	}

	// Display MySQL install page on backend
	if this.Wrap.IsBackend && !this.Wrap.ConfMysqlExists {
		utils.SystemRenderTemplate(this.Wrap.W, assets.TmplCpMySql, nil)
		return true
	}

	// Read MySQL settings file
	mc, err := utils.MySqlConfigRead(this.Wrap.DConfig + string(os.PathSeparator) + "mysql.json")
	if err != nil {
		utils.SystemErrorPageEngine(this.Wrap.W, err)
		return true
	}

	// Connect to MySQL server
	db, err := sql.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
	if err != nil {
		utils.SystemErrorPageEngine(this.Wrap.W, err)
		return true
	}
	this.Wrap.DB = db
	defer db.Close()
	err = db.Ping()

	// Check if MySQL server alive
	if err != nil {
		utils.SystemErrorPageEngine(this.Wrap.W, err)
		return true
	}

	// Front-end or back-end here...
	if this.Wrap.R.URL.Path == "/" {
		this.Wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		this.Wrap.W.Header().Set("Content-Type", "text/html")

		counter := this.Wrap.S.GetInt("counter", 0)
		// this.Wrap.LogAccess(fmt.Sprintf("Counter value was: %d", counter))

		this.Wrap.W.Write([]byte(`Logic -> (` + fmt.Sprintf("%d", counter) + `)`))

		counter++
		this.Wrap.S.SetInt("counter", counter)
		// this.Wrap.LogAccess(fmt.Sprintf("Counter value now: %d", counter))

		return true
	}

	return false
}
