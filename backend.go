package main

import (
	"golang-fave/engine/wrapper"

	Templates "golang-fave/engine/wrapper/resources/templates"
)

func handleBackEnd(e *wrapper.Wrapper) bool {
	return e.TmplBackEnd(Templates.CpLogin, nil)
}
