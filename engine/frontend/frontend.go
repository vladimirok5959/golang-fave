package frontend

import (
	"database/sql"

	"golang-fave/engine/wrapper"
)

type Frontend struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
}

func New(wrapper *wrapper.Wrapper, db *sql.DB) *Frontend {
	return &Frontend{wrapper, db}
}

func (this *Frontend) Run() bool {
	if this.wrapper.R.URL.Path == "/" {
		return this.wrapper.TmplFrontEnd("index", nil)
	}

	return false
}
