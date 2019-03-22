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

		err := wrap.DBTrans(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM blog_cats FOR UPDATE;"); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM blog_cat_post_rel WHERE category_id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec("SELECT @ml := lft, @mr := rgt FROM blog_cats WHERE id = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM blog_cats WHERE id = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET lft = lft - 1, rgt = rgt - 1 WHERE lft > @ml AND rgt < @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET lft = lft - 2 WHERE lft > @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt - 2 WHERE rgt > @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM blog_cat_post_rel WHERE category_id = ?;", pf_id); err != nil {
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
