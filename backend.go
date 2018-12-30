package main

import (
	"golang-fave/engine/wrapper"
	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

func handlerBackEnd(wrapper *wrapper.Wrapper) bool {
	// MySQL config page
	if !utils.IsMySqlConfigExists(wrapper.DirVHostHome) {
		return wrapper.TmplBackEnd(templates.CpMySQL, nil)
	}

	// Login page
	return wrapper.TmplBackEnd(templates.CpLogin, nil)
}
