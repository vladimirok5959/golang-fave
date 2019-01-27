package frontend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
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
	user    *utils.MySql_user
	urls    *[]string
}

func New(wrapper *wrapper.Wrapper, db *sql.DB, url_args *[]string) *Frontend {
	return &Frontend{wrapper, db, nil, url_args}
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
