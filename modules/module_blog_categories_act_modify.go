package modules

import (
	// "golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) blog_ActionCategoryAdd(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	return wrap.DB.Transaction(func(tx *wrapper.Tx) error {
		// Block rows
		if _, err := tx.Exec("SELECT id FROM blog_cats FOR UPDATE;"); err != nil {
			return err
		}

		// Process
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
		return wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Process
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

	// TODO: Fix parent change

	// Parent is changed, move category to new parent
	return wrap.DB.Transaction(func(tx *wrapper.Tx) error {
		// Block all rows
		if _, err := tx.Exec("SELECT id FROM blog_cats FOR UPDATE;"); err != nil {
			return err
		}

		wrap.LogError("Start update!")
		wrap.LogError("--------------------------------")

		var parentL int
		var parentR int
		if err := tx.QueryRow(`SELECT lft, rgt FROM blog_cats WHERE id = ?;`, pf_parent).Scan(&parentL, &parentR); err != nil {
			return err
		}

		var targetL int
		var targetR int
		if err := tx.QueryRow(`SELECT lft, rgt FROM blog_cats WHERE id = ?;`, pf_id).Scan(&targetL, &targetR); err != nil {
			return err
		}

		wrap.LogError("parentL = %d, parentR = %d", parentL, parentR)
		wrap.LogError("targetL = %d, targetR = %d", targetL, targetR)

		// Select data
		rows, err := tx.Query("SELECT id, lft, rgt FROM blog_cats WHERE lft >= ? and rgt <= ? ORDER BY lft ASC", targetL, targetR)
		if err != nil {
			return err
		}
		defer rows.Close()
		var rows_id []int
		var rows_lft []int
		var rows_rgt []int
		for rows.Next() {
			var row_id int
			var row_lft int
			var row_rgt int
			if err := rows.Scan(&row_id, &row_lft, &row_rgt); err == nil {
				rows_id = append(rows_id, row_id)
				rows_lft = append(rows_lft, row_lft)
				rows_rgt = append(rows_rgt, row_rgt)
			}
		}

		wrap.LogError("rows_id = %v", rows_id)
		wrap.LogError("rows_lft = %v", rows_lft)
		wrap.LogError("rows_rgt = %v", rows_rgt)

		if targetL > parentR {
			// From right to left

			// Shift
			step := targetR - targetL + 1
			if _, err := tx.Exec("UPDATE blog_cats SET lft = lft + ? WHERE lft > ? and lft < ?;", step, parentR, targetL); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + ? WHERE rgt > ? and rgt < ?;", step, parentR, targetL); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt + ? WHERE id = ?;", step, pf_parent); err != nil {
				return err
			}

			// Update target rows
			for i, _ := range rows_id {
				new_lft := rows_lft[i] - (targetL - parentR)
				new_rgt := rows_rgt[i] - (targetL - parentR)
				if _, err := tx.Exec("UPDATE blog_cats SET lft = ?, rgt = ? WHERE id = ?;", new_lft, new_rgt, rows_id[i]); err != nil {
					return err
				}
			}
		} else {
			// From left to right

			// Shift
			step := targetR - targetL + 1
			// 8 - 3 + 1 = 6
			if _, err := tx.Exec("UPDATE blog_cats SET lft = lft - ? WHERE lft > ? and lft < ?;", step, targetR, parentR); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE blog_cats SET rgt = rgt - ? WHERE rgt > ? and rgt < ?;", step, targetR, parentR); err != nil {
				return err
			}

			// Update target rows
			for i, _ := range rows_id {
				new_lft := rows_lft[i] + (parentR - targetL - step)
				new_rgt := rows_rgt[i] + (parentR - targetL - step)
				if _, err := tx.Exec("UPDATE blog_cats SET lft = ?, rgt = ? WHERE id = ?;", new_lft, new_rgt, rows_id[i]); err != nil {
					return err
				}
			}
		}

		// Update target cat data
		if _, err := tx.Exec("UPDATE blog_cats SET name = ?, alias = ? WHERE id = ?;", pf_name, pf_alias, pf_id); err != nil {
			return err
		}

		wrap.LogError("--------------------------------")

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
