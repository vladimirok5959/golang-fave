package main

import (
	"net/http"

	"golang-fave/engine/wrapper"
	utils "golang-fave/engine/wrapper/utils"
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

func handleFrontEnd(wrapper *wrapper.Wrapper) bool {
	// Redirect to CP, if MySQL config file is not exists
	if !utils.IsMySqlConfigExists(wrapper.DirVhostHome) {
		(*wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(*wrapper.W, wrapper.R, wrapper.R.URL.Scheme+"://"+wrapper.R.Host+"/cp/", 302)
		return true
	}

	// Else logic here
	if wrapper.R.URL.Path == "/" {
		return wrapper.TmplFrontEnd("index", TmplData{
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
