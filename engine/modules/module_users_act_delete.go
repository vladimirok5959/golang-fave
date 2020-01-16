package modules

import (
	"context"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_UsersDelete() *Action {
	return this.newAction(AInfo{
		Mount:     "users-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_blog_cats WHERE user = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_blog_posts WHERE user = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_pages WHERE user = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "SELECT id FROM fave_users WHERE id = ? and id > 1 FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec(ctx, "UPDATE fave_blog_cats SET user = 1 WHERE user = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "UPDATE fave_blog_posts SET user = 1 WHERE user = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "UPDATE fave_pages SET user = 1 WHERE user = ?;", utils.StrToInt(pf_id)); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "DELETE FROM fave_users WHERE id = ? and id > 1;", utils.StrToInt(pf_id)); err != nil {
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
