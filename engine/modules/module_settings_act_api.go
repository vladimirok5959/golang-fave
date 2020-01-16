package modules

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsApi() *Action {
	return this.newAction(AInfo{
		Mount:     "settings-api",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_xml_enabled := utils.Trim(wrap.R.FormValue("xml-enabled"))
		pf_xml_name := utils.Trim(wrap.R.FormValue("xml-name"))
		pf_xml_company := utils.Trim(wrap.R.FormValue("xml-company"))
		pf_xml_url := utils.Trim(wrap.R.FormValue("xml-url"))

		if pf_xml_enabled == "" {
			pf_xml_enabled = "0"
		}

		(*wrap.Config).API.XML.Enabled = utils.StrToInt(pf_xml_enabled)
		(*wrap.Config).API.XML.Name = pf_xml_name
		(*wrap.Config).API.XML.Company = pf_xml_company
		(*wrap.Config).API.XML.Url = pf_xml_url

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.RecreateProductXmlFile()

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
