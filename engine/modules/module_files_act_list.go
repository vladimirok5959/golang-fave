package modules

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"golang-fave/engine/assets"
	"golang-fave/engine/builder"
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_FilesList() *Action {
	return this.newAction(AInfo{
		Mount:     "files-list",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_path := utils.SafeFilePath(utils.Trim(wrap.R.FormValue("path")))

		// Set path
		wrap.Write(`fave.FilesManagerSetPath('` + pf_path + `');`)

		// Render table
		start_dir := strings.Join([]string{wrap.DHtdocs, "public"}, string(os.PathSeparator)) + pf_path + "*"

		str_dirs := ""
		str_files := ""

		nothing := true

		if files, err := filepath.Glob(start_dir); err == nil {
			for _, file := range files {
				file_name := file
				i := strings.LastIndex(file_name, string(os.PathSeparator))
				if i != -1 {
					file_name = file_name[i+1:]
				}

				if utils.IsDir(file) {
					if nothing {
						nothing = false
					}
					actions := builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   "/public" + pf_path + file_name + "/",
							Hint:   "View",
							Target: "_blank",
						},
						{
							Icon:    assets.SysSvgIconRemove,
							Href:    "javascript:fave.FilesManagerRemoveFolder(\\'" + pf_path + file_name + "\\',\\'Are you sure want to delete folder?\\');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
					str_dirs += `<tr class="dir"><td class="col_name"><a href="javascript:fave.FilesManagerLoadData(\'` + pf_path + file_name + `/` + `\');"><span class="text-dotted">` + html.EscapeString(file_name) + `</span></a></td><td class="col_type"><b>DIR</b></td><td class="col_action">` + actions + `</td></tr>`
				} else {
					if nothing {
						nothing = false
					}
					actions := builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   "/public" + pf_path + file_name,
							Hint:   "View",
							Target: "_blank",
						},
						{
							Icon:    assets.SysSvgIconRemove,
							Href:    "javascript:fave.FilesManagerRemoveFile(\\'" + pf_path + file_name + "\\',\\'Are you sure want to delete file?\\');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})

					size_value := float64(utils.GetFileSize(file))
					size_text := "Bytes"

					if size_value > 1024 {
						size_value = size_value / 1024
						size_text = "Kb"
					}
					if size_value > 1024 {
						size_value = size_value / 1024
						size_text = "Mb"
					}
					if size_value > 1024 {
						size_value = size_value / 1024
						size_text = "Gb"
					}

					str_size := fmt.Sprintf("%5.2f %s", size_value, size_text)
					str_files += `<tr class="file"><td class="col_name"><span class="text-dotted">` + html.EscapeString(file_name) + `</span></td><td class="col_type">` + str_size + `</td><td class="col_action">` + actions + `</td></tr>`
				}
			}
		}

		if pf_path != "/" {
			str_dirs = `<tr class="dir"><td class="col_name"><a href="javascript:fave.FilesManagerLoadDataUp(\'` + pf_path + `\');">..</a></td><td class="col_type">&nbsp;</td><td class="col_action">&nbsp;</td></tr>` + str_dirs
		}

		str_nothing := ``
		if nothing {
			str_nothing = `<tr><td colspan="50">No data</td></tr>`
		}

		table := `<table class="table data-table table-striped table-bordered table-hover table_fm_files"><thead><tr><th class="col_name">File name</th><th class="col_type">Size</th><th class="col_action">Action</th></tr></thead><tbody>` + str_dirs + str_files + str_nothing + `</tbody></table>`
		wrap.Write(`$('#sys-modal-files-manager .dialog-data').html('` + table + `');`)

		// Enable buttons
		wrap.Write(`fave.FilesManagerEnableDisableButtons(false);`)
	})
}
