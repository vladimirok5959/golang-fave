package backend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
)

type Backend struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
}

func New(wrapper *wrapper.Wrapper, db *sql.DB) *Backend {
	return &Backend{wrapper, db}
}

func (this *Backend) Run() bool {
	// TODO:
	// Check if any user exists
	// If not - display form to create first user

	// Login page
	return this.wrapper.TmplBackEnd(templates.CpLogin, nil)
}
