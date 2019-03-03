package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type FERData struct {
	Wrap    *wrapper.Wrapper
	DataRow interface{}
	Is404   bool

	bufferUser *utils.MySql_user
}

func New(wrap *wrapper.Wrapper, drow interface{}, is404 bool) *FERData {
	fer := &FERData{
		Wrap:    wrap,
		DataRow: drow,
		Is404:   is404,
	}
	return fer.init()
}

func (this *FERData) init() *FERData {
	if this.DataRow != nil {
		if this.Wrap.CurrModule == "index" {
			if this.DataRow.(*utils.MySql_page).A_meta_title == "" {
				this.DataRow.(*utils.MySql_page).A_meta_title = this.DataRow.(*utils.MySql_page).A_name
			}
		}
	}
	return this
}
