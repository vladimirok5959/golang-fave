package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_UsersDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "users-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM blog_cats WHERE user = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM blog_posts WHERE user = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM pages WHERE user = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM users WHERE id = ? and id > 1 FOR UPDATE;", pf_id); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec("UPDATE blog_cats SET user = 1 WHERE user = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_posts SET user = 1 WHERE user = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE pages SET user = 1 WHERE user = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM users WHERE id = ? and id > 1;", pf_id); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}