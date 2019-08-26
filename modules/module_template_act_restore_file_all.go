package modules

import (
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/assets/template"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_TemplateRestoreFileAll() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "template-restore-file-all",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_file := wrap.R.FormValue("file")

		if pf_file != "all" {
			wrap.MsgError(`Inner system error`)
			return
		}

		for file, data := range ThemeFiles.AllData {
			if err := ioutil.WriteFile(wrap.DTemplate+string(os.PathSeparator)+file, data, 0664); err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
