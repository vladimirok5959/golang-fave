package frontend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"
)

// --- Demo
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

// --------

type Frontend struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
}

func New(wrapper *wrapper.Wrapper, db *sql.DB) *Frontend {
	return &Frontend{wrapper, db}
}

func (this *Frontend) Run() bool {
	// --- Demo
	if this.wrapper.R.URL.Path == "/" {
		return this.wrapper.TmplFrontEnd("index", TmplData{
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
	// --------

	return false
}
