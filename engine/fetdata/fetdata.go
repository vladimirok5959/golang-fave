package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	wrap    *wrapper.Wrapper
	dataRow interface{}
	is404   bool

	bufferUser *utils.MySql_user
}

func New(wrap *wrapper.Wrapper, drow interface{}, is404 bool) *FERData {
	fer := &FERData{
		wrap:    wrap,
		dataRow: drow,
		is404:   is404,
	}
	return fer.init()
}

func (this *FERData) init() *FERData {
	if this.dataRow != nil {
		if this.wrap.CurrModule == "index" {
			if this.dataRow.(*utils.MySql_page).A_meta_title == "" {
				this.dataRow.(*utils.MySql_page).A_meta_title = this.dataRow.(*utils.MySql_page).A_name
			}
		}
	}
	return this
}

func (this *FERData) Is404() bool {
	return this.is404
}
