package wrapper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"golang-fave/consts"
	"golang-fave/engine/mysqlpool"
	"golang-fave/logger"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Wrapper struct {
	l *logger.Logger
	W http.ResponseWriter
	R *http.Request
	S *session.Session

	Host     string
	Port     string
	CurrHost string

	DConfig   string
	DHtdocs   string
	DLogs     string
	DTemplate string
	DTmp      string

	IsBackend       bool
	ConfMysqlExists bool
	UrlArgs         []string
	CurrModule      string
	CurrSubModule   string
	MSPool          *mysqlpool.MySqlPool

	DB   *sql.DB
	User *utils.MySql_user
}

func New(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, host, port, chost, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string, mp *mysqlpool.MySqlPool) *Wrapper {
	return &Wrapper{
		l:             l,
		W:             w,
		R:             r,
		S:             s,
		Host:          host,
		Port:          port,
		CurrHost:      chost,
		DConfig:       dirConfig,
		DHtdocs:       dirHtdocs,
		DLogs:         dirLogs,
		DTemplate:     dirTemplate,
		DTmp:          dirTmp,
		UrlArgs:       []string{},
		CurrModule:    "",
		CurrSubModule: "",
		MSPool:        mp,
	}
}

func (this *Wrapper) LogAccess(msg string) {
	this.l.Log(msg, this.R, false)
}

func (this *Wrapper) LogError(msg string) {
	this.l.Log(msg, this.R, true)
}

func (this *Wrapper) dbReconnect() error {
	if !utils.IsMySqlConfigExists(this.DConfig + string(os.PathSeparator) + "mysql.json") {
		return errors.New("can't read database configuration file")
	}
	mc, err := utils.MySqlConfigRead(this.DConfig + string(os.PathSeparator) + "mysql.json")
	if err != nil {
		return err
	}
	this.DB, err = sql.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
	if err != nil {
		return err
	}
	this.MSPool.Set(this.CurrHost, this.DB)
	return nil
}

func (this *Wrapper) UseDatabase() error {
	this.DB = this.MSPool.Get(this.CurrHost)
	if this.DB == nil {
		if err := this.dbReconnect(); err != nil {
			return err
		}
	}

	if err := this.DB.Ping(); err != nil {
		this.DB.Close()
		if err := this.dbReconnect(); err != nil {
			return err
		}
		if err := this.DB.Ping(); err != nil {
			this.DB.Close()
			return err
		}
	}

	// Here we are connected
	this.DB.SetConnMaxLifetime(time.Minute * 30)
	this.DB.SetMaxIdleConns(2)
	this.DB.SetMaxOpenConns(2)

	return nil
}

func (this *Wrapper) LoadSessionUser() bool {
	if this.S.GetInt("UserId", 0) <= 0 {
		return false
	}
	if this.DB == nil {
		return false
	}
	user := &utils.MySql_user{}
	err := this.DB.QueryRow(`
		SELECT
			id,
			first_name,
			last_name,
			email,
			password,
			admin,
			active
		FROM
			users
		WHERE
			id = ?
		LIMIT 1;`,
		this.S.GetInt("UserId", 0),
	).Scan(
		&user.A_id,
		&user.A_first_name,
		&user.A_last_name,
		&user.A_email,
		&user.A_password,
		&user.A_admin,
		&user.A_active,
	)
	if err != nil {
		return false
	}
	if user.A_id != this.S.GetInt("UserId", 0) {
		return false
	}
	this.User = user
	return true
}

func (this *Wrapper) Write(data string) {
	this.W.Write([]byte(data))
}

func (this *Wrapper) MsgSuccess(msg string) {
	this.Write(fmt.Sprintf(
		`fave.ShowMsgSuccess('Success!', '%s', false);`,
		utils.JavaScriptVarValue(msg)))
}

func (this *Wrapper) MsgError(msg string) {
	this.Write(fmt.Sprintf(
		`fave.ShowMsgError('Error!', '%s', true);`,
		utils.JavaScriptVarValue(msg)))
}

func (this *Wrapper) RenderToString(tcont []byte, data interface{}) string {
	tmpl, err := template.New("template").Parse(string(tcont))
	if err != nil {
		return err.Error()
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err.Error()
	}
	return tpl.String()
}

func (this *Wrapper) RenderFrontEnd(tname string, data interface{}, status int) {
	tmpl, err := template.ParseFiles(
		this.DTemplate+string(os.PathSeparator)+tname+".html",
		this.DTemplate+string(os.PathSeparator)+"header.html",
		this.DTemplate+string(os.PathSeparator)+"sidebar-left.html",
		this.DTemplate+string(os.PathSeparator)+"sidebar-right.html",
		this.DTemplate+string(os.PathSeparator)+"footer.html",
	)
	if err != nil {
		utils.SystemErrorPageTemplate(this.W, err)
		return
	}
	this.W.WriteHeader(status)
	this.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	this.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, consts.TmplData{
		System: utils.GetTmplSystemData(),
		Data:   data,
	})
	if err != nil {
		utils.SystemErrorPageTemplate(this.W, err)
		return
	}
	this.W.Write(tpl.Bytes())
}

func (this *Wrapper) RenderBackEnd(tcont []byte, data interface{}) {
	tmpl, err := template.New("template").Parse(string(tcont))
	if err != nil {
		utils.SystemErrorPageEngine(this.W, err)
		return
	}
	this.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	this.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tpl bytes.Buffer
	err = tmpl.Execute(this.W, consts.TmplData{
		System: utils.GetTmplSystemData(),
		Data:   data,
	})
	if err != nil {
		utils.SystemErrorPageEngine(this.W, err)
		return
	}
	this.W.Write(tpl.Bytes())
}

func (this *Wrapper) DBTrans(queries func(tx *sql.Tx) error) error {
	if queries == nil {
		return errors.New("queries is not set for transaction")
	}

	tx, err := this.DB.Begin()
	if err != nil {
		return err
	}

	err = queries(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
