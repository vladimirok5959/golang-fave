package main

import (
	"golang-fave/engine/wrapper"
	"net/http"
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
	// Redirect to CP, if MySQL config file is not exists
	if !e.IsMySqlConfigExists() {
		(*e.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(*e.W, e.R, e.R.URL.Scheme+"://"+e.R.Host+"/cp/", 302)
		return true
	}

	// Else logic here
	if e.R.URL.Path == "/" {
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
	return false
}
