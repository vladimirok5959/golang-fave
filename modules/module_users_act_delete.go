package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

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

		err := wrap.DBTrans(func(tx *sql.Tx) error {
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
