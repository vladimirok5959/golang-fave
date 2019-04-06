package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_IndexModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "index-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_content := wrap.R.FormValue("content")
		pf_meta_title := wrap.R.FormValue("meta_title")
		pf_meta_keywords := wrap.R.FormValue("meta_keywords")
		pf_meta_description := wrap.R.FormValue("meta_description")
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
			pf_alias = utils.GenerateAlias(pf_name)
		}

		if !utils.IsValidAlias(pf_alias) {
			wrap.MsgError(`Please specify correct page alias`)
			return
		}

		if pf_id == "0" {
			// Add new page
			_, err := wrap.DB.Exec(
				`INSERT INTO pages SET
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
				pf_active,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/';`)
		} else {
			// Update page
			_, err := wrap.DB.Exec(
				`UPDATE pages SET
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
				pf_active,
				utils.StrToInt(pf_id),
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/index/modify/` + pf_id + `/';`)
		}
	})
}
