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
	// Show add user form if no any user in db
	var count int
	err := this.db.QueryRow("SELECT COUNT(*) FROM `users`;").Scan(&count)
	if this.wrapper.EngineErrMsgOnError(err) {
		return true
	}
	if count <= 0 {
		return this.wrapper.TmplBackEnd(templates.CpFirstUser, nil)
	}

	// Login page
	if this.wrapper.Session.GetIntDef("UserId", 0) <= 0 {
		return this.wrapper.TmplBackEnd(templates.CpLogin, nil)
	}

	(*this.wrapper.W).Write([]byte(`Admin panel here...`))
	return true
}
