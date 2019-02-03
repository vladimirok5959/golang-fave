package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	h http.Handler
	w io.Writer
}

func (this handler) logRequest(w *writer, r *http.Request) {
	fmt.Fprintln(os.Stdout, strings.Join([]string{
		r.Host,
		r.RemoteAddr,
		"-",
		"[" + w.s.Format(time.RFC3339) + "]",
		`"` + r.Method,
		r.RequestURI,
		r.Proto + `"`,
		strconv.Itoa(w.status),
		strconv.Itoa(w.size),
	}, " "))
}

func (this handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrt := &writer{w: w, s: time.Now()}
	this.h.ServeHTTP(wrt, r)
	this.logRequest(wrt, r)
}
