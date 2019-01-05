package main

import (
	"database/sql"
	"net/http"
	"strings"
	//"log"

	"golang-fave/engine/backend"
	"golang-fave/engine/frontend"
	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

func handlerPage(wrapper *wrapper.Wrapper) bool {
	mysql_conf_exists := utils.IsMySqlConfigExists(wrapper.DirVHostHome)

	is_front_end := true
	if wrapper.R.URL.Path == "/cp" || strings.HasPrefix(wrapper.R.URL.Path, "/cp/") {
		is_front_end = false
	}

	if is_front_end {
		// Front-end
		if !mysql_conf_exists {
			(*wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			http.Redirect(*wrapper.W, wrapper.R, wrapper.R.URL.Scheme+"://"+wrapper.R.Host+"/cp/", 302)
			return true
		}
	} else {
		// Back-end
		if !mysql_conf_exists {
			return wrapper.TmplBackEnd(templates.CpMySQL, nil)
		}
	}

	// Connect to database or show error
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

	// Parse url
	/*
		url_buff := wrapper.R.URL.Path
		if len(url_buff) >= 1 && url_buff[:1] == "/" {
			url_buff = url_buff[1:]
		}
		if len(url_buff) >= 1 && url_buff[len(url_buff)-1:] == "/" {
			url_buff = url_buff[:len(url_buff)-1]
		}

		log.Printf("###############")
		log.Printf("(%s)", url_buff)
	*/

	/*
		url_args := utils.UrlToArray(wrapper.R.URL.Path)
		log.Printf("############### (%d)", len(url_args))
		for key, value := range url_args {
			log.Printf(">>> (%d) -> (%s)", key, value)
		}
	*/

	// log.Printf("###############")

	url_args := utils.UrlToArray(wrapper.R.URL.Path)

	// Run WebSite or CP
	if is_front_end {
		// Front-end
		return frontend.New(wrapper, db, &url_args).Run()
	} else {
		// Back-end
		return backend.New(wrapper, db, &url_args).Run()
	}
}
