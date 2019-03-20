package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) blog_ActionCategoryAdd(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	return wrap.DBTrans(func(tx *sql.Tx) error {
		if _, err := tx.Exec("SELECT @mr := rgt FROM blog_cats WHERE id = ?;", pf_parent); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + 2 WHERE rgt > @mr;"); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET lft = lft + 2 WHERE lft > @mr;"); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + 2 WHERE id = ?;", pf_parent); err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT INTO blog_cats (id, user, name, alias, lft, rgt) VALUES (NULL, ?, ?, ?, @mr, @mr + 1);", wrap.User.A_id, pf_name, pf_alias); err != nil {
			return err
		}
		return nil
	})
}

func (this *Modules) blog_ActionCategoryUpdate(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	parentId := this.blog_GetCategoryParentId(wrap, utils.StrToInt(pf_id))

	if utils.StrToInt(pf_parent) == parentId {
		// If parent not changed, just update category data
		return wrap.DBTrans(func(tx *sql.Tx) error {
			if _, err := tx.Exec(`
				UPDATE blog_cats SET
					name = ?,
					alias = ?
				WHERE
					id > 1 AND
					id = ?
				;`,
				pf_name,
				pf_alias,
				pf_id,
			); err != nil {
				return err
			}
			return nil
		})
	}

	// Parent is changed, move category to new parent
	return wrap.DBTrans(func(tx *sql.Tx) error {
		// Shift
		if _, err := tx.Exec("SELECT @ml := lft, @mr := rgt FROM blog_cats WHERE id = ?;", pf_id); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET lft = 0, rgt = 0 WHERE id = ?;", pf_id); err != nil {
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

		// Update
		if _, err := tx.Exec("SELECT @mr := rgt FROM blog_cats WHERE id = ?;", pf_parent); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + 2 WHERE rgt > @mr;"); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET lft = lft + 2 WHERE lft > @mr;"); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + 2 WHERE id = ?;", pf_parent); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE blog_cats SET name = ?, alias = ?, lft = @mr, rgt = @mr + 1 WHERE id = ?;", pf_name, pf_alias, pf_id); err != nil {
			return err
		}

		return nil
	})
}

func (this *Modules) RegisterAction_BlogCategoriesModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-categories-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_parent := wrap.R.FormValue("parent")

		if !utils.IsNumeric(pf_id) || !utils.IsNumeric(pf_parent) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify category name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateSingleAlias(pf_name)
		}

		if !utils.IsValidSingleAlias(pf_alias) {
			wrap.MsgError(`Please specify correct category alias`)
			return
		}

		// Set root category as default
		if pf_parent == "0" {
			pf_parent = "1"
		} else {
			// Check if parent category exists
			var parentId int
			err := wrap.DB.QueryRow(`
				SELECT
					id
				FROM
					blog_cats
				WHERE
					id > 1 AND
					id = ?
				LIMIT 1;`,
				pf_parent,
			).Scan(&parentId)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		if pf_id == "0" {
			if err := this.blog_ActionCategoryAdd(wrap, pf_id, pf_name, pf_alias, pf_parent); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/blog/categories/';`)
		} else {
			if err := this.blog_ActionCategoryUpdate(wrap, pf_id, pf_name, pf_alias, pf_parent); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/blog/categories-modify/` + pf_id + `/';`)
		}
	})
}
