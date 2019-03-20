package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_BlogCategoriesDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-categories-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) || utils.StrToInt(pf_id) <= 1 {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Start transaction
		tx, err := wrap.DB.Begin()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Update and delete target category
		if _, err = tx.Exec("SELECT @ml := lft, @mr := rgt FROM blog_cats WHERE id = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("DELETE FROM blog_cats WHERE id = ?;", pf_id); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("UPDATE blog_cats SET lft = lft - 1, rgt = rgt - 1 WHERE lft > @ml AND rgt < @mr;"); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("UPDATE blog_cats SET lft = lft - 2 WHERE lft > @mr;"); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("UPDATE blog_cats SET rgt = rgt - 2 WHERE rgt > @mr;"); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec("DELETE FROM blog_cat_post_rel WHERE category_id = ?;", pf_id); err != nil {
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
