package modules

import (
	"io/ioutil"
	"os"

	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsRobotsTxt() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-robots-txt",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_content := wrap.R.FormValue("content")

		// Save robots.txt content
		err := ioutil.WriteFile(wrap.DTemplate+string(os.PathSeparator)+"robots.txt", []byte(pf_content), 0664)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
