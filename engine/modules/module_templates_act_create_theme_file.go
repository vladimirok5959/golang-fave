package modules

import (
	"io/ioutil"
	"os"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_TemplatesCreateThemeFile() *Action {
	return this.newAction(AInfo{
		Mount:     "templates-create-theme-file",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_name := utils.Trim(wrap.R.FormValue("name"))
		pf_content := wrap.R.FormValue("content")

		if pf_name == "" {
			wrap.MsgError(`Please specify file name`)
			return
		}

		// Check normal file name here

		template_file := wrap.DTemplate + string(os.PathSeparator) + pf_name + ".html"
		if utils.IsFileExists(template_file) {
			wrap.MsgError(`File is already exists`)
			return
		}

		// Save content to file
		err := ioutil.WriteFile(template_file, []byte(pf_content), 0664)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Redirect to created file in editor

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
