package engine

import (
	"fmt"
	"net/http"

	"golang-fave/engine/wrapper"
	"golang-fave/logger"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

type Engine struct {
	Wrap *wrapper.Wrapper
	// Database
	// Actions
	// Front-end or Back-end
}

func Response(l *logger.Logger, w http.ResponseWriter, r *http.Request, s *session.Session, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp string) bool {
	wrap := wrapper.New(l, w, r, s, host, port, dirConfig, dirHtdocs, dirLogs, dirTemplate, dirTmp)
	eng := &Engine{Wrap: wrap}

	return eng.Process()
}

func (this *Engine) Process() bool {
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
