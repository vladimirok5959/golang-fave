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

	"golang-fave/engine/basket"
	"golang-fave/engine/cblocks"
	"golang-fave/engine/config"
	"golang-fave/engine/consts"
	"golang-fave/engine/logger"
	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"

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
	ShopBasket      *basket.Basket
	Config          *config.Config

	DB   *sqlw.DB
	User *utils.MySql_user

	ShopAllCurrencies *map[int]utils.MySql_shop_currency
}

func New(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, c *cblocks.CacheBlocks, host, port, chost, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string, mp *mysqlpool.MySqlPool, sb *basket.Basket) *Wrapper {
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
		ShopBasket:    sb,
		Config:        conf,
	}
}

func (this *Wrapper) LogAccess(msg string, vars ...interface{}) {
	this.l.Log(msg, this.R, false, vars...)
}

func (this *Wrapper) LogError(msg string, vars ...interface{}) {
	this.l.Log(msg, this.R, true, vars...)
}

func (this *Wrapper) LogCpError(err *error) *error {
	if *err != nil {
		this.LogError("%s", *err)
	}
	return err
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

	// Max 60 minutes and max 8 connection per host
	this.DB.SetConnMaxLifetime(time.Minute * 60)
	this.DB.SetMaxIdleConns(8)
	this.DB.SetMaxOpenConns(8)

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
	if err := this.DB.Ping(this.R.Context()); err != nil {
		this.DB.Close()
		if err := this.dbReconnect(); err != nil {
			return err
		} else {
			if err := this.DB.Ping(this.R.Context()); err != nil {
				this.DB.Close()
				this.MSPool.Del(this.CurrHost)
				return err
			}
		}
	}
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
	err := this.DB.QueryRow(
		this.R.Context(),
		`SELECT
			id,
			first_name,
			last_name,
			email,
			password,
			admin,
			active
		FROM
			fave_users
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
	if *this.LogCpError(&err) != nil {
		return false
	}
	if user.A_id != this.S.GetInt("UserId", 0) {
		return false
	}
	this.User = user
	return true
}

func (this *Wrapper) Write(data string) {
	// TODO: select context and don't write
	this.W.Write([]byte(data))
}

func (this *Wrapper) MsgSuccess(msg string) {
	// TODO: select context and don't write
	this.Write(fmt.Sprintf(
		`fave.ShowMsgSuccess('Success!', '%s', false);`,
		utils.JavaScriptVarValue(msg)))
}

func (this *Wrapper) MsgError(msg string) {
	// TODO: select context and don't write
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
	templates := []string{
		this.DTemplate + string(os.PathSeparator) + tname + ".html",
		this.DTemplate + string(os.PathSeparator) + "header.html",
		this.DTemplate + string(os.PathSeparator) + "sidebar-left.html",
		this.DTemplate + string(os.PathSeparator) + "sidebar-right.html",
		this.DTemplate + string(os.PathSeparator) + "footer.html",
	}
	tmpl, err := template.New(tname + ".html").Funcs(utils.TemplateAdditionalFuncs()).ParseFiles(templates...)
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

func (this *Wrapper) SendEmailUsual(email, subject, message string) error {
	if (*this.Config).SMTP.Host == "" || (*this.Config).SMTP.Login == "" || (*this.Config).SMTP.Password == "" {
		return errors.New("SMTP server is not configured")
	}

	err := utils.SMTPSend(
		(*this.Config).SMTP.Host,
		utils.IntToStr((*this.Config).SMTP.Port),
		(*this.Config).SMTP.Login,
		(*this.Config).SMTP.Password,
		subject,
		message,
		[]string{email},
	)

	status := 1
	emessage := ""
	if err != nil {
		status = 0
		emessage = err.Error()
	}

	if _, err := this.DB.Exec(
		this.R.Context(),
		`INSERT INTO fave_notify_mail SET
			id = NULL,
			email = ?,
			subject = ?,
			message = ?,
			error = ?,
			datetime = ?,
			status = ?
		;`,
		email,
		subject,
		message,
		emessage,
		utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
		status,
	); err != nil {
		this.LogCpError(&err)
	}

	return err
}

func (this *Wrapper) SendEmailFast(email, subject, message string) error {
	status := 2
	emessage := ""

	if (*this.Config).SMTP.Host == "" || (*this.Config).SMTP.Login == "" || (*this.Config).SMTP.Password == "" {
		status = 0
		emessage = "SMTP server is not configured"
	}

	if _, err := this.DB.Exec(
		this.R.Context(),
		`INSERT INTO fave_notify_mail SET
			id = NULL,
			email = ?,
			subject = ?,
			message = ?,
			error = ?,
			datetime = ?,
			status = ?
		;`,
		email,
		subject,
		message,
		emessage,
		utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
		status,
	); err != nil {
		return err
	}
	return nil
}

func (this *Wrapper) SendEmailTemplated(email, subject, tname string, data interface{}) error {
	tmpl, err := template.New(tname + ".html").Funcs(utils.TemplateAdditionalFuncs()).ParseFiles(
		this.DTemplate + string(os.PathSeparator) + tname + ".html",
	)
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return err
	}

	return this.SendEmailFast(email, subject, string(tpl.Bytes()))
}

func (this *Wrapper) GetSessionId() string {
	cookie, err := this.R.Cookie("session")
	if err == nil && len(cookie.Value) == 40 {
		return cookie.Value
	}
	return ""
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
		// TODO: select context and don't do that
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				return errors.New(fmt.Sprintf("[upload delete] Thumbnail file (%s) delete error: %s", file, err.Error()))
			}
		}
	}
	return this.RecreateProductImgFiles()
}

func (this *Wrapper) ShopGetAllCurrencies() *map[int]utils.MySql_shop_currency {
	if this.ShopAllCurrencies == nil {
		this.ShopAllCurrencies = &map[int]utils.MySql_shop_currency{}
		if rows, err := this.DB.Query(
			this.R.Context(),
			`SELECT
				id,
				name,
				coefficient,
				code,
				symbol
			FROM
				fave_shop_currencies
			ORDER BY
				id ASC
			;`,
		); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := utils.MySql_shop_currency{}
				if err = rows.Scan(
					&row.A_id,
					&row.A_name,
					&row.A_coefficient,
					&row.A_code,
					&row.A_symbol,
				); err == nil {
					(*this.ShopAllCurrencies)[row.A_id] = row
				}
			}
		}
	}
	return this.ShopAllCurrencies
}

func (this *Wrapper) ShopGetCurrentCurrency() *utils.MySql_shop_currency {
	currency_id := 1
	if cookie, err := this.R.Cookie("currency"); err == nil {
		currency_id = utils.StrToInt(cookie.Value)
	}
	if _, ok := (*this.ShopGetAllCurrencies())[currency_id]; ok != true {
		currency_id = 1
	}
	if p, ok := (*this.ShopGetAllCurrencies())[currency_id]; ok == true {
		return &p
	}
	return nil
}

func (this *Wrapper) IsSystemMountedTemplateFile(filename string) bool {
	return utils.InArrayString([]string{
		"404.html",
		"blog-category.html",
		"blog-post.html",
		"blog.html",
		"cached-block-1.html",
		"cached-block-2.html",
		"cached-block-3.html",
		"cached-block-4.html",
		"cached-block-5.html",
		"email-new-order-admin.html",
		"email-new-order-user.html",
		"footer.html",
		"header.html",
		"index.html",
		"maintenance.html",
		"page.html",
		"robots.txt",
		"scripts.js",
		"shop-category.html",
		"shop-product.html",
		"shop.html",
		"sidebar-left.html",
		"sidebar-right.html",
		"styles.css",
	}, filename)
}
