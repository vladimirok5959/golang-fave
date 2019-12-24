package modules

import (
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/engine/assets/template"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_TemplateRestoreFile() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "template-restore-file",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_file := utils.Trim(wrap.R.FormValue("file"))

		if pf_file == "" {
			wrap.MsgError(`Please specify file name`)
			return
		}

		if _, ok := ThemeFiles.AllData[pf_file]; !ok {
			wrap.MsgError(`Template file is not defined in engine`)
			return
		}

		// Restore file content
		err := ioutil.WriteFile(wrap.DTemplate+string(os.PathSeparator)+pf_file, ThemeFiles.AllData[pf_file], 0664)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
