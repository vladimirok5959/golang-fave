package main

import (
	"golang-fave/engine/wrapper"

	Templates "golang-fave/engine/wrapper/resources/templates"
)

func handleBackEnd(e *wrapper.Wrapper) bool {
	// MySQL config page
	if !e.IsMySqlConfigExists() {
		return e.TmplBackEnd(Templates.CpMySQL, nil)
	}

	// Login page
	return e.TmplBackEnd(Templates.CpLogin, nil)
}
