package modules

import (
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_BlogModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_content := wrap.R.FormValue("content")
		pf_active := wrap.R.FormValue("active")

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify page name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateSingleAlias(pf_name)
		}

		if !utils.IsValidSingleAlias(pf_alias) {
			wrap.MsgError(`Please specify correct post alias`)
			return
		}

		if pf_id == "0" {
			// Start transaction
			tx, err := wrap.DB.Begin()
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}

			// Insert row
			res, err := tx.Exec(
				`INSERT INTO blog_posts SET
					user = ?,
					name = ?,
					alias = ?,
					content = ?,
					datetime = ?,
					active = ?
				;`,
				wrap.User.A_id,
				pf_name,
				pf_alias,
				pf_content,
				utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
				pf_active,
			)
			if err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}

			lastID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}

			// Insert post and categories relations
			catids := utils.GetPostArrayInt("cats[]", wrap.R)
			if len(catids) > 0 {
				var catsCount int
				err = tx.QueryRow(`
					SELECT
						COUNT(*)
					FROM
						blog_cats
					WHERE
						id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
					;`,
				).Scan(
					&catsCount,
				)
				if err != nil {
					tx.Rollback()
					wrap.MsgError(err.Error())
					return
				}
				if len(catids) != catsCount {
					tx.Rollback()
					wrap.MsgError(`Inner system error`)
					return
				}
				var balkInsertArr []string
				for _, el := range catids {
					balkInsertArr = append(balkInsertArr, `(NULL,`+utils.Int64ToStr(lastID)+`,`+utils.IntToStr(el)+`)`)
				}
				if _, err = tx.Exec(
					`INSERT INTO blog_cat_post_rel (id,post_id,category_id) VALUES ` + strings.Join(balkInsertArr, ",") + `;`,
				); err != nil {
					tx.Rollback()
					wrap.MsgError(err.Error())
					return
				}
			}

			// Commit all changes
			err = tx.Commit()
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.Write(`window.location='/cp/blog/';`)
		} else {
			// Start transaction
			tx, err := wrap.DB.Begin()
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}

			// Update row
			if _, err = tx.Exec(
				`UPDATE blog_posts SET
					name = ?,
					alias = ?,
					content = ?,
					active = ?
				WHERE
					id = ?
				;`,
				pf_name,
				pf_alias,
				pf_content,
				pf_active,
				utils.StrToInt(pf_id),
			); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}

			// Delete post and categories relations
			if _, err = tx.Exec("DELETE FROM blog_cat_post_rel WHERE post_id = ?;", pf_id); err != nil {
				tx.Rollback()
				wrap.MsgError(err.Error())
				return
			}

			// Insert post and categories relations
			catids := utils.GetPostArrayInt("cats[]", wrap.R)
			if len(catids) > 0 {
				var catsCount int
				err = tx.QueryRow(`
					SELECT
						COUNT(*)
					FROM
						blog_cats
					WHERE
						id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
					;`,
				).Scan(
					&catsCount,
				)
				if err != nil {
					tx.Rollback()
					wrap.MsgError(err.Error())
					return
				}
				if len(catids) != catsCount {
					tx.Rollback()
					wrap.MsgError(`Inner system error`)
					return
				}
				var balkInsertArr []string
				for _, el := range catids {
					balkInsertArr = append(balkInsertArr, `(NULL,`+pf_id+`,`+utils.IntToStr(el)+`)`)
				}
				if _, err = tx.Exec(
					`INSERT INTO blog_cat_post_rel (id,post_id,category_id) VALUES ` + strings.Join(balkInsertArr, ",") + `;`,
				); err != nil {
					tx.Rollback()
					wrap.MsgError(err.Error())
					return
				}
			}

			// Commit all changes
			err = tx.Commit()
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}

			wrap.Write(`window.location='/cp/blog/modify/` + pf_id + `/';`)
		}
	})
}
