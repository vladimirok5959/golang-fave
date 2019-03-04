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
	msg := fmt.Sprint(strings.Join([]string{
		r.Host,
		r.RemoteAddr,
		"-",
		"[" + w.s.Format(time.RFC3339) + "]",
		`"` + r.Method,
		r.RequestURI,
		r.Proto + `"`,
		strconv.Itoa(w.status),
		strconv.Itoa(w.size),
		fmt.Sprintf("%.3f ms", time.Now().Sub(w.s).Seconds()),
	}, " "))

	// Do not wait
	go func() {
		select {
		case this.c <- logMsg{r.Host, msg, w.status >= 400}:
			return
		case <-time.After(1 * time.Second):
			fmt.Println("Logger error, log channel is overflowed (2)")
			return
		}
	}()
}

func (this handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrt := &writer{w: w, s: time.Now()}
	this.h.ServeHTTP(wrt, r)
	this.log(wrt, r)
}
