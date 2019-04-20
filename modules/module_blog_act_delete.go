package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_BlogDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM blog_posts WHERE id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM blog_cat_post_rel WHERE post_id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}

			// Delete target post with category connection data
			if _, err := tx.Exec("DELETE FROM blog_cat_post_rel WHERE post_id = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM blog_posts WHERE id = ?;", pf_id); err != nil {
				return err
			}
			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
