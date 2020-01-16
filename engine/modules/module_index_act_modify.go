package modules

import (
	"context"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_IndexModify() *Action {
	return this.newAction(AInfo{
		Mount:     "index-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_name := utils.Trim(wrap.R.FormValue("name"))
		pf_alias := utils.Trim(wrap.R.FormValue("alias"))
		pf_content := utils.Trim(wrap.R.FormValue("content"))
		pf_meta_title := utils.Trim(wrap.R.FormValue("meta_title"))
		pf_meta_keywords := utils.Trim(wrap.R.FormValue("meta_keywords"))
		pf_meta_description := utils.Trim(wrap.R.FormValue("meta_description"))
		pf_active := utils.Trim(wrap.R.FormValue("active"))

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
			pf_alias = utils.GenerateAlias(pf_name)
		}

		if !utils.IsValidAlias(pf_alias) {
			wrap.MsgError(`Please specify correct page alias`)
			return
		}

		if pf_id == "0" {
			// Add new page
			var lastID int64 = 0
			if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
				res, err := tx.Exec(
					ctx,
					`INSERT INTO fave_pages SET
						user = ?,
						name = ?,
						alias = ?,
						content = ?,
						meta_title = ?,
						meta_keywords = ?,
						meta_description = ?,
						datetime = ?,
						active = ?
					;`,
					wrap.User.A_id,
					pf_name,
					pf_alias,
					pf_content,
					pf_meta_title,
					pf_meta_keywords,
					pf_meta_description,
					utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
					utils.StrToInt(pf_active),
				)
				if err != nil {
					return err
				}
				// Get inserted post id
				lastID, err = res.LastInsertId()
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/index/modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			// Update page
			if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
				_, err := tx.Exec(
					ctx,
					`UPDATE fave_pages SET
						name = ?,
						alias = ?,
						content = ?,
						meta_title = ?,
						meta_keywords = ?,
						meta_description = ?,
						active = ?
					WHERE
						id = ?
					;`,
					pf_name,
					pf_alias,
					pf_content,
					pf_meta_title,
					pf_meta_keywords,
					pf_meta_description,
					utils.StrToInt(pf_active),
					utils.StrToInt(pf_id),
				)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/index/modify/` + pf_id + `/';`)
		}
	})
}
