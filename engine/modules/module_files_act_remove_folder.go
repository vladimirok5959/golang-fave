package modules

import (
	"os"
	"strings"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_FilesRemoveFolder() *Action {
	return this.newAction(AInfo{
		Mount:     "files-remove-folder",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_file := utils.SafeFilePath(utils.Trim(wrap.R.FormValue("file")))

		file := strings.Join([]string{wrap.DHtdocs, "public"}, string(os.PathSeparator)) + pf_file
		if err := os.RemoveAll(file); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		path := "/"
		i := strings.LastIndex(pf_file, string(os.PathSeparator))
		if i != -1 {
			path = pf_file[:i+1]
		}

		// Set path
		wrap.Write(`fave.FilesManagerSetPath('` + path + `');`)

		// Refresh table
		wrap.Write(`fave.FilesManagerLoadData('` + path + `');`)
	})
}
