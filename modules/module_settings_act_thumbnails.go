package modules

import (
	"strconv"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_SettingsThumbnails() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-thumbnails",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_shop_thumbnail_w_1 := wrap.R.FormValue("shop-thumbnail-w-1")
		pf_shop_thumbnail_h_1 := wrap.R.FormValue("shop-thumbnail-h-1")
		pf_shop_thumbnail_r_1 := wrap.R.FormValue("shop-thumbnail-r-1")

		pf_shop_thumbnail_w_2 := wrap.R.FormValue("shop-thumbnail-w-2")
		pf_shop_thumbnail_h_2 := wrap.R.FormValue("shop-thumbnail-h-2")
		pf_shop_thumbnail_r_2 := wrap.R.FormValue("shop-thumbnail-r-2")

		pf_shop_thumbnail_w_3 := wrap.R.FormValue("shop-thumbnail-w-3")
		pf_shop_thumbnail_h_3 := wrap.R.FormValue("shop-thumbnail-h-3")
		pf_shop_thumbnail_r_3 := wrap.R.FormValue("shop-thumbnail-r-3")

		pf_shop_thumbnail_w_full := wrap.R.FormValue("shop-thumbnail-w-full")
		pf_shop_thumbnail_h_full := wrap.R.FormValue("shop-thumbnail-h-full")
		pf_shop_thumbnail_r_full := wrap.R.FormValue("shop-thumbnail-r-full")

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_1); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_1); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_r_1); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_2); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_2); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_r_2); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_3); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_3); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_r_3); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_full); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_full); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_r_full); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_shop_thumbnail_w_1 := utils.StrToInt(pf_shop_thumbnail_w_1)
		pfi_shop_thumbnail_h_1 := utils.StrToInt(pf_shop_thumbnail_h_1)
		pfi_shop_thumbnail_r_1 := utils.StrToInt(pf_shop_thumbnail_r_1)

		pfi_shop_thumbnail_w_2 := utils.StrToInt(pf_shop_thumbnail_w_2)
		pfi_shop_thumbnail_h_2 := utils.StrToInt(pf_shop_thumbnail_h_2)
		pfi_shop_thumbnail_r_2 := utils.StrToInt(pf_shop_thumbnail_r_2)

		pfi_shop_thumbnail_w_3 := utils.StrToInt(pf_shop_thumbnail_w_3)
		pfi_shop_thumbnail_h_3 := utils.StrToInt(pf_shop_thumbnail_h_3)
		pfi_shop_thumbnail_r_3 := utils.StrToInt(pf_shop_thumbnail_r_3)

		pfi_shop_thumbnail_w_full := utils.StrToInt(pf_shop_thumbnail_w_full)
		pfi_shop_thumbnail_h_full := utils.StrToInt(pf_shop_thumbnail_h_full)
		pfi_shop_thumbnail_r_full := utils.StrToInt(pf_shop_thumbnail_r_full)

		// Correct some values
		if pfi_shop_thumbnail_w_1 < 10 {
			pfi_shop_thumbnail_w_1 = 10
		}
		if pfi_shop_thumbnail_h_1 > 1000 {
			pfi_shop_thumbnail_h_1 = 1000
		}
		if pfi_shop_thumbnail_r_1 > 1 {
			pfi_shop_thumbnail_r_1 = 1
		}
		if pfi_shop_thumbnail_r_1 < 0 {
			pfi_shop_thumbnail_r_1 = 0
		}

		if pfi_shop_thumbnail_w_2 < 10 {
			pfi_shop_thumbnail_w_2 = 10
		}
		if pfi_shop_thumbnail_h_2 > 1000 {
			pfi_shop_thumbnail_h_2 = 1000
		}
		if pfi_shop_thumbnail_r_2 > 1 {
			pfi_shop_thumbnail_r_2 = 1
		}
		if pfi_shop_thumbnail_r_2 < 0 {
			pfi_shop_thumbnail_r_2 = 0
		}

		if pfi_shop_thumbnail_w_3 < 10 {
			pfi_shop_thumbnail_w_3 = 10
		}
		if pfi_shop_thumbnail_h_3 > 1000 {
			pfi_shop_thumbnail_h_3 = 1000
		}
		if pfi_shop_thumbnail_r_3 > 1 {
			pfi_shop_thumbnail_r_3 = 1
		}
		if pfi_shop_thumbnail_r_3 < 0 {
			pfi_shop_thumbnail_r_3 = 0
		}

		if pfi_shop_thumbnail_w_full < 10 {
			pfi_shop_thumbnail_w_full = 10
		}
		if pfi_shop_thumbnail_h_full > 1000 {
			pfi_shop_thumbnail_h_full = 1000
		}
		if pfi_shop_thumbnail_r_full > 1 {
			pfi_shop_thumbnail_r_full = 1
		}
		if pfi_shop_thumbnail_r_full < 0 {
			pfi_shop_thumbnail_r_full = 0
		}

		is_changed_tb1 := false
		is_changed_tb2 := false
		is_changed_tb3 := false
		is_changed_tbf := false

		if (*wrap.Config).Shop.Thumbnails.Thumbnail1[0] != pfi_shop_thumbnail_w_1 || (*wrap.Config).Shop.Thumbnails.Thumbnail1[1] != pfi_shop_thumbnail_h_1 || (*wrap.Config).Shop.Thumbnails.Thumbnail1[2] != pfi_shop_thumbnail_r_1 {
			is_changed_tb1 = true
		}
		if (*wrap.Config).Shop.Thumbnails.Thumbnail2[0] != pfi_shop_thumbnail_w_2 || (*wrap.Config).Shop.Thumbnails.Thumbnail2[1] != pfi_shop_thumbnail_h_2 || (*wrap.Config).Shop.Thumbnails.Thumbnail2[2] != pfi_shop_thumbnail_r_2 {
			is_changed_tb2 = true
		}
		if (*wrap.Config).Shop.Thumbnails.Thumbnail3[0] != pfi_shop_thumbnail_w_3 || (*wrap.Config).Shop.Thumbnails.Thumbnail3[1] != pfi_shop_thumbnail_h_3 || (*wrap.Config).Shop.Thumbnails.Thumbnail3[2] != pfi_shop_thumbnail_r_3 {
			is_changed_tb3 = true
		}
		if (*wrap.Config).Shop.Thumbnails.ThumbnailFull[0] != pfi_shop_thumbnail_w_full || (*wrap.Config).Shop.Thumbnails.ThumbnailFull[1] != pfi_shop_thumbnail_h_full || (*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] != pfi_shop_thumbnail_r_full {
			is_changed_tbf = true
		}

		(*wrap.Config).Shop.Thumbnails.Thumbnail1[0] = pfi_shop_thumbnail_w_1
		(*wrap.Config).Shop.Thumbnails.Thumbnail1[1] = pfi_shop_thumbnail_h_1
		(*wrap.Config).Shop.Thumbnails.Thumbnail1[2] = pfi_shop_thumbnail_r_1

		(*wrap.Config).Shop.Thumbnails.Thumbnail2[0] = pfi_shop_thumbnail_w_2
		(*wrap.Config).Shop.Thumbnails.Thumbnail2[1] = pfi_shop_thumbnail_h_2
		(*wrap.Config).Shop.Thumbnails.Thumbnail2[2] = pfi_shop_thumbnail_r_2

		(*wrap.Config).Shop.Thumbnails.Thumbnail3[0] = pfi_shop_thumbnail_w_3
		(*wrap.Config).Shop.Thumbnails.Thumbnail3[1] = pfi_shop_thumbnail_h_3
		(*wrap.Config).Shop.Thumbnails.Thumbnail3[2] = pfi_shop_thumbnail_r_3

		(*wrap.Config).Shop.Thumbnails.ThumbnailFull[0] = pfi_shop_thumbnail_w_full
		(*wrap.Config).Shop.Thumbnails.ThumbnailFull[1] = pfi_shop_thumbnail_h_full
		(*wrap.Config).Shop.Thumbnails.ThumbnailFull[2] = pfi_shop_thumbnail_r_full

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reset products images cache
		if is_changed_tb1 || is_changed_tb2 || is_changed_tb3 || is_changed_tbf {
			if is_changed_tb1 && is_changed_tb2 && is_changed_tb3 && is_changed_tbf {
				if err := wrap.RemoveProductImageThumbnails("*", "thumb-*"); err != nil {
					wrap.MsgError(err.Error())
					return
				}
			} else {
				if is_changed_tb1 {
					if err := wrap.RemoveProductImageThumbnails("*", "thumb-1-*"); err != nil {
						wrap.MsgError(err.Error())
						return
					}
				}
				if is_changed_tb2 {
					if err := wrap.RemoveProductImageThumbnails("*", "thumb-2-*"); err != nil {
						wrap.MsgError(err.Error())
						return
					}
				}
				if is_changed_tb3 {
					if err := wrap.RemoveProductImageThumbnails("*", "thumb-3-*"); err != nil {
						wrap.MsgError(err.Error())
						return
					}
				}
				if is_changed_tbf {
					if err := wrap.RemoveProductImageThumbnails("*", "thumb-full-*"); err != nil {
						wrap.MsgError(err.Error())
						return
					}
				}
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
