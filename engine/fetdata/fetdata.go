package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	Wrap    *wrapper.Wrapper
	DataRow interface{}

	bufferUser *utils.MySql_user
}

func New(wrap *wrapper.Wrapper, drow interface{}) *FERData {
	fer := &FERData{
		Wrap:    wrap,
		DataRow: drow,
	}
	return fer.init()
}

func (this *FERData) init() *FERData {
	if this.Wrap.CurrModule == "index" {
		if this.DataRow.(*utils.MySql_page).A_meta_title == "" {
			this.DataRow.(*utils.MySql_page).A_meta_title = this.DataRow.(*utils.MySql_page).A_name
		}
	}
	return this
}
