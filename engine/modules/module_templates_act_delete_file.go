package modules

import (
	"context"
	"os"
	"strings"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_TemplatesDeleteThemeFile() *Action {
	return this.newAction(AInfo{
		Mount:     "templates-delete-file",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_file := utils.Trim(wrap.R.FormValue("file"))

		if pf_file == "" {
			wrap.MsgError(`Please specify file name`)
			return
		}

		if wrap.IsSystemMountedTemplateFile(pf_file) {
			wrap.MsgError(`You can't delete system mounted template`)
			return
		}

		template_file := wrap.DTemplate + string(os.PathSeparator) + pf_file
		if !utils.IsFileExists(template_file) {
			wrap.MsgError(`File is not exists`)
			return
		}

		if err := os.Remove(template_file); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Update pages
		tmpl_name := pf_file
		if i := strings.LastIndex(tmpl_name, "."); i > -1 {
			tmpl_name = tmpl_name[:i]
		}

		if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			_, err := tx.Exec(
				ctx,
				`UPDATE fave_pages SET
					template = ?
				WHERE
					template = ?
				;`,
				"page",
				tmpl_name,
			)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		wrap.Write(`window.location='/cp/templates/';`)
	})
}
