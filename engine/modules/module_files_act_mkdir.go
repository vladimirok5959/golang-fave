package modules

import (
	"os"
	"strings"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_FilesMkdir() *Action {
	return this.newAction(AInfo{
		Mount:     "files-mkdir",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_path := utils.Trim(wrap.R.FormValue("path"))
		pf_name := utils.Trim(wrap.R.FormValue("name"))

		if pf_path == "" {
			wrap.MsgError(`Please specify folder path`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify folder name`)
			return
		}

		dirname := utils.SafeFilePath(pf_path + pf_name)
		target := strings.Join([]string{wrap.DHtdocs, "public"}, string(os.PathSeparator)) + dirname
		if err := os.Mkdir(target, 0755); err != nil {
			emsg := err.Error()
			i := strings.Index(emsg, ":")
			if i != -1 {
				emsg = emsg[i+1:]
			}
			wrap.MsgError(emsg)
			return
		}

		path := "/"
		i := strings.LastIndex(dirname, string(os.PathSeparator))
		if i != -1 {
			path = dirname[:i+1]
		}

		// Set path
		wrap.Write(`$('#sys-modal-files-manager .dialog-path span').html('` + path + `');`)

		// Refresh table
		wrap.Write(`fave.FilesManagerLoadData('` + path + `');`)
	})
}
