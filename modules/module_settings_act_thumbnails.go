package modules

import (
	"os"
	"path/filepath"
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

		pf_shop_thumbnail_w_2 := wrap.R.FormValue("shop-thumbnail-w-2")
		pf_shop_thumbnail_h_2 := wrap.R.FormValue("shop-thumbnail-h-2")

		pf_shop_thumbnail_w_3 := wrap.R.FormValue("shop-thumbnail-w-3")
		pf_shop_thumbnail_h_3 := wrap.R.FormValue("shop-thumbnail-h-3")

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_1); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_1); err != nil {
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

		if _, err := strconv.Atoi(pf_shop_thumbnail_w_3); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}
		if _, err := strconv.Atoi(pf_shop_thumbnail_h_3); err != nil {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_shop_thumbnail_w_1 := utils.StrToInt(pf_shop_thumbnail_w_1)
		pfi_shop_thumbnail_h_1 := utils.StrToInt(pf_shop_thumbnail_h_1)

		pfi_shop_thumbnail_w_2 := utils.StrToInt(pf_shop_thumbnail_w_2)
		pfi_shop_thumbnail_h_2 := utils.StrToInt(pf_shop_thumbnail_h_2)

		pfi_shop_thumbnail_w_3 := utils.StrToInt(pf_shop_thumbnail_w_3)
		pfi_shop_thumbnail_h_3 := utils.StrToInt(pf_shop_thumbnail_h_3)

		// Correct some values
		if pfi_shop_thumbnail_w_1 < 10 {
			pfi_shop_thumbnail_w_1 = 10
		}
		if pfi_shop_thumbnail_h_1 > 1000 {
			pfi_shop_thumbnail_h_1 = 1000
		}

		if pfi_shop_thumbnail_w_2 < 10 {
			pfi_shop_thumbnail_w_2 = 10
		}
		if pfi_shop_thumbnail_h_2 > 1000 {
			pfi_shop_thumbnail_h_2 = 1000
		}

		if pfi_shop_thumbnail_w_3 < 10 {
			pfi_shop_thumbnail_w_3 = 10
		}
		if pfi_shop_thumbnail_h_3 > 1000 {
			pfi_shop_thumbnail_h_3 = 1000
		}

		(*wrap.Config).Shop.Thumbnails.Thumbnail1[0] = pfi_shop_thumbnail_w_1
		(*wrap.Config).Shop.Thumbnails.Thumbnail1[1] = pfi_shop_thumbnail_h_1

		(*wrap.Config).Shop.Thumbnails.Thumbnail2[0] = pfi_shop_thumbnail_w_2
		(*wrap.Config).Shop.Thumbnails.Thumbnail2[1] = pfi_shop_thumbnail_h_2

		(*wrap.Config).Shop.Thumbnails.Thumbnail3[0] = pfi_shop_thumbnail_w_3
		(*wrap.Config).Shop.Thumbnails.Thumbnail3[1] = pfi_shop_thumbnail_h_3

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reset products images cache
		pattern := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + "*" + string(os.PathSeparator) + "thumb-*"
		if files, err := filepath.Glob(pattern); err == nil {
			for _, file := range files {
				if err := os.Remove(file); err != nil {
					wrap.LogError("Thumbnail file (%s) delete error: %s", file, err.Error())
				}
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
