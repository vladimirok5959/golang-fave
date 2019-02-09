package wrapper

import (
	"database/sql"
	"net/http"

	"golang-fave/logger"

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

	DB *sql.DB
}

func New(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) *Wrapper {
	return &Wrapper{
		l:         l,
		W:         w,
		R:         r,
		S:         s,
		Host:      host,
		Port:      port,
		DConfig:   dirConfig,
		DHtdocs:   dirHtdocs,
		DLogs:     dirLogs,
		DTemplate: dirTemplate,
		DTmp:      dirTmp,
	}
}

func (this *Wrapper) LogAccess(msg string) {
	this.l.Log(msg, this.R, false)
}

func (this *Wrapper) LogError(msg string) {
	this.l.Log(msg, this.R, true)
}
