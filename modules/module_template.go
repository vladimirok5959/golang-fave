package modules

import (
	"html"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang-fave/engine/assets"
	"golang-fave/engine/builder"
	"golang-fave/engine/consts"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) template_GetThemeFiles(wrap *wrapper.Wrapper) []string {
	var result []string
	files, err := ioutil.ReadDir(wrap.DTemplate)
	if err == nil {
		for _, file := range files {
			if len(file.Name()) > 0 && file.Name()[0] == '.' {
				continue
			}
			if len(file.Name()) > 0 && strings.ToLower(file.Name()) == "robots.txt" {
				continue
			}
			result = append(result, file.Name())
		}
	}
	return result
}

func (this *Modules) RegisterModule_Template() *Module {
	return this.newModule(MInfo{
		WantDB: false,
		Mount:  "template",
		Name:   "Template",
		Order:  802,
		System: true,
		Icon:   assets.SysSvgIconView,
		Sub: &[]MISub{
			{Mount: "default", Name: "Editor", Show: true, Icon: assets.SysSvgIconEdit},
			{Mount: "restore", Name: "Restore", Show: true, Icon: assets.SysSvgIconRestore},
		},
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Editor"},
			})

			files := this.template_GetThemeFiles(wrap)
			if len(files) > 0 {
				selected_file, _ := url.QueryUnescape(wrap.R.URL.Query().Get("file"))
				if !(selected_file != "" && utils.InArrayString(files, selected_file)) {
					selected_file = files[0]
				}

				list_of_files := ``
				for _, file := range files {
					selected := ""
					if file == selected_file {
						selected = " selected"
					}
					list_of_files += `<option value="` + html.EscapeString(file) +
						`"` + selected + `>` + html.EscapeString(file) + `</option>`
				}

				fcont := []byte(``)
				fcont, _ = ioutil.ReadFile(wrap.DTemplate + string(os.PathSeparator) + selected_file)

				fext := filepath.Ext(selected_file)
				if len(fext) > 2 {
					fext = fext[1:]
				}

				content += builder.DataForm(wrap, []builder.DataFormField{
					{
						Kind:  builder.DFKHidden,
						Name:  "action",
						Value: "template-edit-theme-file",
					},
					{
						Kind:    builder.DFKText,
						Caption: "Theme file",
						Name:    "file",
						Value:   "0",
						CallBack: func(field *builder.DataFormField) string {
							return `<div class="form-group n1">` +
								`<div class="row">` +
								`<div class="col-12">` +
								`<div style="position:relative;">` +
								`<button type="button" class="btn btn-success" onclick="return fave.ActionRestoreThemeFile('template-restore-file','` + selected_file + `','Are you sure want to restore theme file?');" style="position:absolute;right:0;">Restore</button>` +
								`<select class="form-control ignore-lost-data" id="lbl_file" name="file" onchange="setTimeout(function(){$('#lbl_file').val('` + selected_file + `')},500);document.location='/cp/` + wrap.CurrModule + `/?file='+encodeURI(this.value);">` +
								list_of_files +
								`</select>` +
								`</div>` +
								`</div>` +
								`</div>` +
								`</div>`
						},
					},
					{
						Kind: builder.DFKText,
						CallBack: func(field *builder.DataFormField) string {
							return `<div class="form-group last"><div class="row"><div class="col-12"><textarea class="form-control tmpl-editor" name="content" data-emode="` + fext + `" placeholder="" autocomplete="off">` + html.EscapeString(string(fcont)) + `</textarea></div></div></div>`
						},
					},
					{
						Kind: builder.DFKSubmit,
						CallBack: func(field *builder.DataFormField) string {
							return `<div class="row d-lg-none"><div class="col-12"><div class="pt-3"><button type="submit" class="btn btn-primary" data-target="add-edit-button">Save</button></div></div></div>`
						},
					},
				})

				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			} else {
				content += `<div class="sys-messages">
					<div class="alert alert-warning" role="alert">
						<strong>Error!</strong> No any file found in theme folder
					</div>
				</div>`
			}
		} else if wrap.CurrSubModule == "restore" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Restore"},
			})

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind: builder.DFKText,
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group last"><div class="row"><div class="col-12"><div class="alert alert-danger" style="margin:0;"><strong>WARNING!</strong> This action will restore current theme files to original, you will lost you theme changes! Think twice before run this action! If you still want to do this, please press <b>Restore</b> red button!</div></div></div></div>`
					},
				},
				{
					Kind: builder.DFKSubmit,
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="row d-lg-none"><div class="col-12"><div class="pt-3"><button type="button" class="btn btn-danger" onclick="return fave.ActionRestoreThemeFile('template-restore-file-all','all','WARNING! Are you sure want to restore all theme files?');">Restore</button></div></div></div>`
					},
				},
			})

			sidebar += `<button class="btn btn-danger btn-sidebar" onclick="return fave.ActionRestoreThemeFile('template-restore-file-all','all','WARNING! Are you sure want to restore all theme files?');" id="add-edit-button">Restore</button>`
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
