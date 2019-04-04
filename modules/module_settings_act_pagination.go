package modules

import (
	"strconv"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_SettingsPagination() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-pagination",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_blog_index := wrap.R.FormValue("blog-index")
		pf_blog_category := wrap.R.FormValue("blog-category")

		if _, err := strconv.Atoi(pf_blog_index); err != nil {
			wrap.MsgError(`Blog posts count per page on main page must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_blog_category); err != nil {
			wrap.MsgError(`Blog posts count per page on category page must be integer number`)
			return
		}

		pfi_blog_index := utils.StrToInt(pf_blog_index)
		pfi_blog_category := utils.StrToInt(pf_blog_category)

		// Correct some values
		if pfi_blog_index < 0 {
			pfi_blog_index = 1
		}
		if pfi_blog_index > 100 {
			pfi_blog_index = 100
		}

		if pfi_blog_category < 0 {
			pfi_blog_category = 1
		}
		if pfi_blog_category > 100 {
			pfi_blog_category = 100
		}

		(*wrap.Config).Blog.Pagination.Index = pfi_blog_index
		(*wrap.Config).Blog.Pagination.Category = pfi_blog_category

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
