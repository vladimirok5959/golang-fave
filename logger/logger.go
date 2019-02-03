package logger

import (
	"net/http"
	"os"
)

func New(h http.Handler) http.Handler {
	return handler{
		h: h,
		w: os.Stdout,
	}
}
