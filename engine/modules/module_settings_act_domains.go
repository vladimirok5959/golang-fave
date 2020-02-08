package modules

import (
	"io/ioutil"
	"os"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsDomains() *Action {
	return this.newAction(AInfo{
		Mount:     "settings-domains",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_content := wrap.R.FormValue("content")

		// Save content
		err := ioutil.WriteFile(wrap.DConfig+string(os.PathSeparator)+".domains", []byte(pf_content), 0664)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
