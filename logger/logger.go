package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang-fave/consts"
	"golang-fave/utils"
)

type logMsg struct {
	host    string
	message string
	isError bool
}

type Logger struct {
	wwwDir string
	cdata  chan logMsg
	cclose chan bool
}

func (this *Logger) console(msg *logMsg) {
	if consts.ParamDebug {
		if !msg.isError {
			if consts.IS_WIN {
				fmt.Fprintln(os.Stdout, "[ACCESS] "+msg.message)
			} else {
				fmt.Fprintln(os.Stdout, "\033[0;32m[ACCESS] "+msg.message+"\033[0m")
			}
		} else {
			if consts.IS_WIN {
				fmt.Fprintln(os.Stdout, "[ERROR] "+msg.message)
			} else {
				fmt.Fprintln(os.Stdout, "\033[0;31m[ERROR] "+msg.message+"\033[0m")
			}
		}
		return
	}

	if !msg.isError {
		fmt.Fprintln(os.Stdout, msg.message)
	} else {
		fmt.Fprintln(os.Stderr, msg.message)
	}
}

func (this *Logger) write(msg *logMsg) {
	// Ignore file if debug
	if consts.ParamDebug {
		this.console(msg)
		return
	}

	// Ignore file if host not set
	if msg.host == "" {
		this.console(msg)
		return
	}

	// Ignore file if www dir is not exists
	if !utils.IsDirExists(this.wwwDir) {
		this.console(msg)
		return
	}

	// Extract host
	host, _ := utils.ExtractHostPort(msg.host, false)
	logs_dir := this.wwwDir + string(os.PathSeparator) + host + string(os.PathSeparator) + "logs"

	// Try use localhost folder for logs
	if !utils.IsDirExists(logs_dir) {
		logs_dir = this.wwwDir + string(os.PathSeparator) + "localhost" + string(os.PathSeparator) + "logs"
	}

	// Ignore file if logs dir is not exists
	if !utils.IsDirExists(logs_dir) {
		this.console(msg)
		return
	}

	// Detect which log file
	log_file := logs_dir + string(os.PathSeparator) + "access.log"
	if msg.isError {
		log_file = logs_dir + string(os.PathSeparator) + "error.log"
	}

	// Try write to file
	f, err := os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		defer f.Close()
		fmt.Fprintln(f, msg.message)
		return
	}

	// By default
	this.console(msg)
}

func New() *Logger {
	// Logs channel
	cdata := make(chan logMsg)

	// Close channel
	cclose := make(chan bool)

	// Init logger pointer
	lg := Logger{cdata: cdata, cclose: cclose}

	// Write log string in background
	go func() {
		for {
			select {
			case msg := <-cdata:
				lg.write(&msg)
			case <-cclose:
				cclose <- true
				return
			}
		}
	}()

	return &lg
}

func (this *Logger) Log(msg string, r *http.Request, isError bool, vars ...interface{}) {
	var host string = ""
	if r != nil {
		host = r.Host
	}

	if len(vars) > 0 {
		msg = fmt.Sprintf(msg, vars...)
	}

	select {
	case this.cdata <- logMsg{host, msg, isError}:
		return
	case <-time.After(1 * time.Second):
		fmt.Println("Logger error, log channel is overflowed (1)")
		return
	}
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
