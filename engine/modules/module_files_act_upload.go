package modules

import (
	"io"
	"os"
	"strings"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_FilesUpload() *Action {
	return this.newAction(AInfo{
		Mount:     "files-upload",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_count := utils.Trim(wrap.R.FormValue("count"))
		pf_path := utils.Trim(wrap.R.FormValue("path"))

		if !utils.IsNumeric(pf_count) {
			wrap.MsgError(`Inner system error`)
			return
		}

		pf_count_int := utils.StrToInt(pf_count)
		if pf_count_int <= 0 {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_path == "" {
			wrap.MsgError(`Please specify files path`)
			return
		}

		for i := 1; i <= pf_count_int; i++ {
			post_field_name := "file_" + utils.IntToStr(i-1)
			if file, handler, err := wrap.R.FormFile(post_field_name); err == nil {
				if handler.Filename != "" {
					filename := utils.SafeFilePath(pf_path + handler.Filename)
					target := strings.Join([]string{wrap.DHtdocs, "public"}, string(os.PathSeparator)) + filename
					if f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0666); err == nil {
						io.Copy(f, file)
						f.Close()
					}
				}
				file.Close()
			}
		}

		path := "/"
		i := strings.LastIndex(pf_path, string(os.PathSeparator))
		if i != -1 {
			path = pf_path[:i+1]
		}

		// Set path
		wrap.Write(`$('#sys-modal-files-manager .dialog-path span').html('` + path + `');`)

		// Refresh table
		wrap.Write(`fave.FilesManagerLoadData('` + path + `');`)
	})
}
