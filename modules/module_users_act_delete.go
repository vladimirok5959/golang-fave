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

		// Start transaction
		tx, err := wrap.DB.Begin()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Update and delete target user
		if _, err = tx.Exec("UPDATE blog_cats SET user = 1 WHERE user = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("UPDATE blog_posts SET user = 1 WHERE user = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("UPDATE pages SET user = 1 WHERE user = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("DELETE FROM users WHERE id = ? and id > 1;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Commit all changes
		err = tx.Commit()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
