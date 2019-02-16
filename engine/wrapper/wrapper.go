package wrapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang-fave/logger"
	"golang-fave/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Wrapper struct {
	l *logger.Logger
	W http.ResponseWriter
	R *http.Request
	S *session.Session

	Host string
	Port string

	DConfig   string
	DHtdocs   string
	DLogs     string
	DTemplate string
	DTmp      string

	IsBackend       bool
	ConfMysqlExists bool
	UrlArgs         []string
	CurrModule      string

	DB   *sql.DB
	User *utils.MySql_user
}

func New(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) *Wrapper {
	return &Wrapper{
		l:          l,
		W:          w,
		R:          r,
		S:          s,
		Host:       host,
		Port:       port,
		DConfig:    dirConfig,
		DHtdocs:    dirHtdocs,
		DLogs:      dirLogs,
		DTemplate:  dirTemplate,
		DTmp:       dirTmp,
		UrlArgs:    []string{},
		CurrModule: "",
	}
}

func (this *Wrapper) LogAccess(msg string) {
	this.l.Log(msg, this.R, false)
}

func (this *Wrapper) LogError(msg string) {
	this.l.Log(msg, this.R, true)
}

func (this *Wrapper) UseDatabase() error {
	if this.DB != nil {
		return errors.New("already connected to database")
	}
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
	err = this.DB.Ping()
	if err != nil {
		this.DB.Close()
		return err
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
	err := this.DB.QueryRow("SELECT `id`, `first_name`, `last_name`, `email`, `password` FROM `users` WHERE `id` = ? LIMIT 1;", this.S.GetInt("UserId", 0)).Scan(
		&user.A_id, &user.A_first_name, &user.A_last_name, &user.A_email, &user.A_password)
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
		`ShowSystemMsgSuccess('Success!', '%s', false);`,
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}

func (this *Wrapper) MsgError(msg string) {
	this.Write(fmt.Sprintf(
		`ShowSystemMsgError('Error!', '%s', true);`,
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}
