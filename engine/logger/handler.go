package logger

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	h http.Handler
	w io.Writer
	c chan logMsg
}

func (this handler) log(w *writer, r *http.Request) {
	rip := r.RemoteAddr
	if r.Header.Get("X-Real-IP") != "" && len(r.Header.Get("X-Real-IP")) <= 25 {
		rip = rip + ", " + strings.TrimSpace(r.Header.Get("X-Real-IP"))
	} else if r.Header.Get("X-Forwarded-For") != "" && len(r.Header.Get("X-Forwarded-For")) <= 25 {
		rip = rip + ", " + strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	}

	uagent := "-"
	if r.Header.Get("User-Agent") != "" && len(r.Header.Get("User-Agent")) <= 256 {
		uagent = strings.TrimSpace(r.Header.Get("User-Agent"))
	}

	msg := fmt.Sprint(strings.Join([]string{
		r.Host,
		"(" + rip + ")",
		"[" + w.s.Format(time.RFC3339) + "]",
		`"` + r.Method,
		r.RequestURI,
		r.Proto + `"`,
		strconv.Itoa(w.status),
		strconv.Itoa(w.size),
		fmt.Sprintf("%.3f ms", time.Now().Sub(w.s).Seconds()),
		`"` + uagent + `"`,
	}, " "))

	select {
	case <-r.Context().Done():
		return
	case this.c <- logMsg{r.Host, msg, w.status >= 400}:
		return
	case <-time.After(1 * time.Second):
		fmt.Printf("Logger, can't write msg (overflow): %s\n", msg)
		return
	}
}

func (this handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrt := &writer{w: w, s: time.Now()}
	this.h.ServeHTTP(wrt, r)
	this.log(wrt, r)
}
