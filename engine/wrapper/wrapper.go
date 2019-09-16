package wrapper

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang-fave/cblocks"
	"golang-fave/consts"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper/config"
	"golang-fave/logger"
	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Tx = sqlw.Tx

var ErrNoRows = sqlw.ErrNoRows

type Wrapper struct {
	l *logger.Logger
	W http.ResponseWriter
	R *http.Request
	S *session.Session
	c *cblocks.CacheBlocks

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
	Config          *config.Config

	DB   *sqlw.DB
	User *utils.MySql_user
}

func New(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, c *cblocks.CacheBlocks, host, port, chost, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string, mp *mysqlpool.MySqlPool) *Wrapper {

	conf := config.ConfigNew()
	if err := conf.ConfigRead(dirConfig + string(os.PathSeparator) + "config.json"); err != nil {
		l.Log("Host config file: %s", r, true, err.Error())
	}

	return &Wrapper{
		l:             l,
		W:             w,
		R:             r,
		S:             s,
		c:             c,
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
		Config:        conf,
	}
}

func (this *Wrapper) LogAccess(msg string, vars ...interface{}) {
	this.l.Log(msg, this.R, false, vars...)
}

func (this *Wrapper) LogError(msg string, vars ...interface{}) {
	this.l.Log(msg, this.R, true, vars...)
}

func (this *Wrapper) dbReconnect() error {
	if !utils.IsMySqlConfigExists(this.DConfig + string(os.PathSeparator) + "mysql.json") {
		return errors.New("can't read database configuration file")
	}
	mc, err := utils.MySqlConfigRead(this.DConfig + string(os.PathSeparator) + "mysql.json")
	if err != nil {
		return err
	}
	this.DB, err = sqlw.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
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

	// Max 60 minutes and max 4 connection per host
	this.DB.SetConnMaxLifetime(time.Minute * 60)
	this.DB.SetMaxIdleConns(4)
	this.DB.SetMaxOpenConns(4)

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
	tmpl, err := template.New(tname+".html").Funcs(utils.TemplateAdditionalFuncs()).ParseFiles(
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
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, consts.TmplData{
		System: utils.GetTmplSystemData("", ""),
		Data:   data,
	})
	if err != nil {
		utils.SystemErrorPageTemplate(this.W, err)
		return
	}
	this.W.WriteHeader(status)
	this.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	this.W.Header().Set("Content-Type", "text/html; charset=utf-8")
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
		System: utils.GetTmplSystemData(this.CurrModule, this.CurrSubModule),
		Data:   data,
	})
	if err != nil {
		utils.SystemErrorPageEngine(this.W, err)
		return
	}
	this.W.Write(tpl.Bytes())
}

func (this *Wrapper) GetCurrentPage(max int) int {
	curr := 1
	page := this.R.URL.Query().Get("p")
	if page != "" {
		if i, err := strconv.Atoi(page); err == nil {
			if i < 1 {
				curr = 1
			} else if i > max {
				curr = max
			} else {
				curr = i
			}
		}
	}
	return curr
}

func (this *Wrapper) ConfigSave() error {
	return this.Config.ConfigWrite(this.DConfig + string(os.PathSeparator) + "config.json")
}

func (this *Wrapper) RecreateProductXmlFile() error {
	trigger := strings.Join([]string{this.DTmp, "trigger.xml.run"}, string(os.PathSeparator))
	if !utils.IsFileExists(trigger) {
		if _, err := os.Create(trigger); err != nil {
			return err
		}
	}
	return nil
}

func (this *Wrapper) RecreateProductImgFiles() error {
	trigger := strings.Join([]string{this.DTmp, "trigger.img.run"}, string(os.PathSeparator))
	if !utils.IsFileExists(trigger) {
		if _, err := os.Create(trigger); err != nil {
			return err
		}
	}
	return nil
}

func (this *Wrapper) RemoveProductImageThumbnails(product_id, filename string) error {
	pattern := this.DHtdocs + string(os.PathSeparator) + strings.Join([]string{"products", "images", product_id, filename}, string(os.PathSeparator))
	if files, err := filepath.Glob(pattern); err != nil {
		return err
	} else {
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				return errors.New(fmt.Sprintf("[upload delete] Thumbnail file (%s) delete error: %s", file, err.Error()))
			}
		}
	}
	return this.RecreateProductImgFiles()
}
