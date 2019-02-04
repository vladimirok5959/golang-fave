package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	wwwDir string
	cdata  chan string
	cclose chan bool
}

func (this *Logger) write(str string) {
	// TODO: write to console or to file
	// If www and host home dir are exists
	fmt.Fprintln(os.Stdout, str)
}

func New() *Logger {
	// Logs channel
	cdata := make(chan string)

	// Close channel
	cclose := make(chan bool)

	// Init logger pointer
	lg := Logger{cdata: cdata, cclose: cclose}

	// Write log string in background
	go func() {
		for {
			select {
			case str := <-cdata:
				lg.write(str)
			case <-cclose:
				cclose <- true
				return
			}
		}
	}()

	return &lg
}

func (this *Logger) Log(str string) {
	// Do not wait
	go func() {
		select {
		case this.cdata <- str:
			return
		case <-time.After(1 * time.Second):
			fmt.Println("Logger error, log channel is overflowed (1)")
			return
		}
	}()
}

func (this *Logger) SetWwwDir(dir string) {
	this.wwwDir = dir
}

func (this *Logger) Handler(h http.Handler) http.Handler {
	return handler{
		h: h,
		w: os.Stdout,
		c: this.cdata,
	}
}

func (this *Logger) Close() {
	this.cclose <- true
	<-this.cclose
}
