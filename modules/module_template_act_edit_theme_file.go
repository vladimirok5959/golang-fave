package modules

import (
	"io/ioutil"
	"os"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_TemplateEditThemeFile() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "template-edit-theme-file",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_file := wrap.R.FormValue("file")
		pf_content := wrap.R.FormValue("content")

		if pf_file == "" {
			wrap.MsgError(`Please specify file name`)
			return
		}

		template_file := wrap.DTemplate + string(os.PathSeparator) + pf_file
		if !utils.IsFileExists(template_file) {
			wrap.MsgError(`File is not exists`)
			return
		}
		if utils.IsDir(template_file) {
			wrap.MsgError(`It's not regular file`)
			return
		}

		files := this.template_GetThemeFiles(wrap)
		if len(files) <= 0 {
			wrap.MsgError(`No any file found in theme folder`)
			return
		}

		if !utils.InArrayString(files, pf_file) {
			wrap.MsgError(`File is not found`)
			return
		}

		// Save content to file
		err := ioutil.WriteFile(template_file, []byte(pf_content), 0664)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
