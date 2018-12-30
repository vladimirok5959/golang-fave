package main

import (
	"database/sql"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

func handlerPage(wrapper *wrapper.Wrapper) bool {
	if !(wrapper.R.URL.Path == "/cp" || strings.HasPrefix(wrapper.R.URL.Path, "/cp/")) {
		return handlerFrontEnd(wrapper)
	} else {
		return handlerBackEnd(wrapper)
	}
}

func handlerFrontEnd(wrapper *wrapper.Wrapper) bool {
	// Redirect to CP, if MySQL config file is not exists
	if !utils.IsMySqlConfigExists(wrapper.DirVHostHome) {
		(*wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(*wrapper.W, wrapper.R, wrapper.R.URL.Scheme+"://"+wrapper.R.Host+"/cp/", 302)
		return true
	}

	// Connect to database

	// Else logic here
	if wrapper.R.URL.Path == "/" {
		return wrapper.TmplFrontEnd("index", nil)
	}

	return false
}

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
