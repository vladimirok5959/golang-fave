package engine

import (
	"fmt"
	"net/http"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Engine struct {
	w http.ResponseWriter
	r *http.Request
	s *session.Session

	host string
	port string

	dConfig   string
	dHtdocs   string
	dLogs     string
	dTemplate string
	dTmp      string
}

func New(w http.ResponseWriter, r *http.Request, s *session.Session, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) *Engine {
	return &Engine{w, r, s, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp}
}

func (this *Engine) Response() bool {
	if this.r.URL.Path == "/" {
		this.w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		this.w.Header().Set("Content-Type", "text/html")

		counter := this.s.GetInt("counter", 0)
		this.w.Write([]byte(`Logic -> (` + fmt.Sprintf("%d", counter) + `)`))

		counter++
		this.s.SetInt("counter", counter)

		return true
	}
	return false
}
