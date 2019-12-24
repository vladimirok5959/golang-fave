package logger

import (
	"net/http"
	"time"
)

type writer struct {
	w http.ResponseWriter
	s time.Time

	status int
	size   int
}

func (this *writer) Header() http.Header {
	return this.w.Header()
}

func (this *writer) Write(bytes []byte) (int, error) {
	if this.status == 0 {
		this.status = http.StatusOK
	}
	size, err := this.w.Write(bytes)
	this.size += size
	return size, err
}

func (this *writer) WriteHeader(status int) {
	this.status = status
	this.w.WriteHeader(status)
}
