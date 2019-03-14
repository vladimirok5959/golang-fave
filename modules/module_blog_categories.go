package modules

import (
	"html"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) blog_GetCategorySelectOptions(wrap *wrapper.Wrapper, id int, parentId int) string {
	result := ``
	rows, err := wrap.DB.Query(
		`SELECT
			node.id,
			node.user,
			node.name,
			node.alias,
			(COUNT(parent.id) - 1) AS depth
		FROM
			blog_cats AS node,
			blog_cats AS parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt AND
			node.id > 1
		GROUP BY
			node.id
		ORDER BY
			node.lft ASC
		;`,
	)
	if err == nil {
		values := make([]string, 5)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		idStr := utils.IntToStr(id)
		parentIdStr := utils.IntToStr(parentId)
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				disabled := ""
				if string(values[0]) == idStr {
					disabled = " disabled"
				}
				selected := ""
				if string(values[0]) == parentIdStr {
					selected = " selected"
				}
				depth := utils.StrToInt(string(values[4])) - 1
				if depth < 0 {
					depth = 0
				}
				sub := strings.Repeat("&mdash; ", depth)
				result += `<option value="` + html.EscapeString(string(values[0])) + `"` + disabled + selected + `>` + sub + html.EscapeString(string(values[2])) + `</option>`
			}
		}
	}
	return result
}

func (this *Modules) blog_GetCategoryParentId(wrap *wrapper.Wrapper, id int) int {
	var parentId int
	_ = wrap.DB.QueryRow(`
		SELECT
			parent.id
		FROM
			blog_cats AS node,
			blog_cats AS parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt AND
			node.id = ? AND
			parent.id <> ?
		ORDER BY
			parent.lft DESC
		LIMIT 1;`,
		id,
		id,
	).Scan(
		&parentId,
	)
	return parentId
}

func (this *Modules) blog_ActionCategoryAdd(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	// Start transaction with table lock
	_, err := wrap.DB.Exec("LOCK TABLE blog_cats WRITE;")
	if err != nil {
		return err
	}
	tx, err := wrap.DB.Begin()
	if err != nil {
		return err
	}

	// Update and insert new category
	if _, err = tx.Exec("SELECT @mr := rgt FROM blog_cats WHERE id = ?;", pf_parent); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("UPDATE blog_cats SET rgt = rgt + 2  WHERE rgt > @mr;"); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("UPDATE blog_cats SET lft = lft + 2  WHERE lft > @mr;"); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("UPDATE blog_cats SET rgt = rgt + 2  WHERE id = ?;", pf_parent); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("INSERT INTO blog_cats (id, user, name, alias, lft, rgt) VALUES (NULL, ?, ?, ?, @mr, @mr + 1);", wrap.User.A_id, pf_name, pf_alias); err != nil {
		tx.Rollback()
		return err
	}

	// Commit all changes and unlock table
	err = tx.Commit()
	if err != nil {
		return err
	}
	_, err = wrap.DB.Exec("UNLOCK TABLES;")
	if err != nil {
		return err
	}

	return nil
}

func (this *Modules) blog_ActionCategoryUpdate(wrap *wrapper.Wrapper, pf_id, pf_name, pf_alias, pf_parent string) error {
	_, err := wrap.DB.Query(
		`UPDATE blog_cats SET
			name = ?,
			alias = ?
		WHERE
			id = ?
		;`,
		pf_name,
		pf_alias,
		pf_id,
	)
	return err
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

func (this *Modules) RegisterAction_BlogCategoriesDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-categories-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Start transaction with table lock
		_, err := wrap.DB.Exec("LOCK TABLE blog_cats WRITE;")
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		tx, err := wrap.DB.Begin()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Update and insert new category
		/*
			if _, err = tx.Exec("SELECT @ml := lft, @mr := rgt, @mw := rgt - lft + 1 FROM blog_cats WHERE id = ?;", pf_id); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}
			if _, err = tx.Exec("DELETE FROM blog_cats WHERE lft BETWEEN @ml AND @mr AND id > 1;"); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}
			if _, err = tx.Exec("UPDATE blog_cats SET rgt = rgt - @mw WHERE rgt > @mr;"); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}
			if _, err = tx.Exec("UPDATE blog_cats SET lft = lft - @mw WHERE lft > @mr;"); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}
		*/

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
