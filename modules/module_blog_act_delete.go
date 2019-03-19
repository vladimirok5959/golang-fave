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

		// Start transaction with table lock
		_, err := wrap.DB.Exec("LOCK TABLES blog_posts WRITE, blog_cat_post_rel WRITE;")
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		tx, err := wrap.DB.Begin()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Delete target post with category connection data
		if _, err = tx.Exec("DELETE FROM blog_posts WHERE id = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("DELETE FROM blog_cat_post_rel WHERE post_id = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Commit all changes and unlock table
		err = tx.Commit()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = wrap.DB.Exec("UNLOCK TABLES;")
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
