package main

import (
	"html/template"

	"golang-fave/engine/wrapper"
)

type TmplMenuItem struct {
	Name string
	Link string
	Active bool
}

type TmplData struct {
	PathSysIcoFav       string
	PathSysCssBootstrap string
	PathSysJsJquery     string
	PathSysJsPopper     string
	PathSysJsBootstrap  string

	MetaTitle       string
	MetaKeywords    string
	MetaDescription string
	MenuItems       []TmplMenuItem
	SomeHtml        template.HTML
}

func handleFrontEnd(e *wrapper.Wrapper) bool {
	tmpl, err := template.ParseFiles(
		e.DirVhostHome+"/template"+"/index.html",
		e.DirVhostHome+"/template"+"/header.html",
		e.DirVhostHome+"/template"+"/sidebar.html",
		e.DirVhostHome+"/template"+"/footer.html",
	)
	if err != nil {
		e.PrintTmplPageError(err)
		return true
	}

	tmpl.Execute(*e.W, TmplData{
		PathSysIcoFav:       e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/fave.ico",
		PathSysCssBootstrap: e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/bootstrap.css",
		PathSysJsJquery:     e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/jquery.js",
		PathSysJsPopper:     e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/popper.js",
		PathSysJsBootstrap:  e.R.URL.Scheme + "://" + e.R.Host + "/assets/sys/bootstrap.js",

		MetaTitle:       "Meta Title",
		MetaKeywords:    "Meta Keywords",
		MetaDescription: "Meta Description",

		MenuItems: []TmplMenuItem{
			{Name: "Home", Link: "/", Active: true},
			{Name: "Item 1", Link: "/#1", Active: false},
			{Name: "Item 2", Link: "/#2", Active: false},
			{Name: "Item 3", Link: "/#3", Active: false},
		},
		SomeHtml: template.HTML("<div class=\"some-class\">DIV</div>"),
	})
	return true
}
