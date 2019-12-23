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
		pf_id := utils.Trim(wrap.R.FormValue("id"))

		if !utils.IsNumeric(pf_id) || utils.StrToInt(pf_id) <= 1 {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(wrap.R.Context(), func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM blog_cats FOR UPDATE;"); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT category_id FROM blog_cat_post_rel WHERE category_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT id FROM blog_posts WHERE category = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Set root category
			if _, err := tx.Exec("UPDATE blog_posts SET category = 1 WHERE category = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec("DELETE FROM blog_cat_post_rel WHERE category_id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT @ml := lft, @mr := rgt FROM blog_cats WHERE id = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM blog_cats WHERE id = ?;", utils.StrToInt(pf_id)); err != nil {
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
			return nil
		})

		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
