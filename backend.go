package main

/*
import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"
	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

func handlerBackEnd(wrapper *wrapper.Wrapper) bool {
	// MySQL config page
	if !utils.IsMySqlConfigExists(wrapper.DirVHostHome) {
		return wrapper.TmplBackEnd(templates.CpMySQL, nil)
	}

	// Connect to database
	mc, err := utils.MySqlConfigRead(wrapper.DirVHostHome)
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}
	db, err := sql.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}
	defer db.Close()
	err = db.Ping()
	if wrapper.EngineErrMsgOnError(err) {
		return true
	}

	// Check if any user exists

	// Login page
	return wrapper.TmplBackEnd(templates.CpLogin, nil)
}
*/
