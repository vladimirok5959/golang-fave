package main

import (
	"golang-fave/engine/wrapper"
)

type MenuItem struct {
	Name   string
	Link   string
	Active bool
}

type TmplData struct {
	MetaTitle       string
	MetaKeywords    string
	MetaDescription string
	MenuItems       []MenuItem
}

func handleFrontEnd(e *wrapper.Wrapper) bool {
	return e.TmplFrontEnd("index", TmplData{
		MetaTitle:       "Meta Title",
		MetaKeywords:    "Meta Keywords",
		MetaDescription: "Meta Description",

		MenuItems: []MenuItem{
			{Name: "Home", Link: "/", Active: true},
			{Name: "Item 1", Link: "/#1", Active: false},
			{Name: "Item 2", Link: "/#2", Active: false},
			{Name: "Item 3", Link: "/#3", Active: false},
		},
	})
}
