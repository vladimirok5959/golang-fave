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
	c chan string
}

func (this handler) log(w *writer, r *http.Request) {
	str := fmt.Sprint(strings.Join([]string{
		r.Host,
		r.RemoteAddr,
		"-",
		"[" + w.s.Format(time.RFC3339) + "]",
		`"` + r.Method,
		r.RequestURI,
		r.Proto + `"`,
		strconv.Itoa(w.status),
		strconv.Itoa(w.size),
		fmt.Sprintf("%.3f ms", time.Now().Sub(w.s).Seconds()/1e6),
	}, " "))

	// Do not wait
	go func() {
		select {
		case this.c <- str:
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
