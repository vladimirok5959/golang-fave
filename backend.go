package main

import (
	"golang-fave/engine/wrapper"

	Templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

func handleBackEnd(e *wrapper.Wrapper) bool {
	// MySQL config page
	if !utils.IsMySqlConfigExists(e.DirVhostHome) {
		return e.TmplBackEnd(Templates.CpMySQL, nil)
	}

	// Login page
	return e.TmplBackEnd(Templates.CpLogin, nil)
}
