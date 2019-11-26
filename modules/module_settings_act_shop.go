package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_SettingsShop() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "settings-shop",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_price_fomat := wrap.R.FormValue("price-fomat")
		pf_price_round := wrap.R.FormValue("price-round")

		pf_require_last_name := wrap.R.FormValue("require-last-name")
		pf_require_first_name := wrap.R.FormValue("require-first-name")
		pf_require_middle_name := wrap.R.FormValue("require-middle-name")
		pf_require_mobile_phone := wrap.R.FormValue("require-mobile-phone")
		pf_require_email_address := wrap.R.FormValue("require-email-address")
		pf_require_delivery := wrap.R.FormValue("require-delivery")
		pf_require_comment := wrap.R.FormValue("require-comment")

		if !utils.IsNumeric(pf_price_fomat) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		if !utils.IsNumeric(pf_price_round) {
			wrap.MsgError(`Must be integer number`)
			return
		}

		pfi_price_fomat := utils.StrToInt(pf_price_fomat)
		pfi_price_round := utils.StrToInt(pf_price_round)

		// Correct values
		if pfi_price_fomat < 0 {
			pfi_price_fomat = 0
		}
		if pfi_price_fomat > 4 {
			pfi_price_fomat = 4
		}

		if pfi_price_round < 0 {
			pfi_price_round = 0
		}
		if pfi_price_round > 2 {
			pfi_price_round = 2
		}

		if pf_require_last_name == "" {
			pf_require_last_name = "0"
		}
		if pf_require_first_name == "" {
			pf_require_first_name = "0"
		}
		if pf_require_middle_name == "" {
			pf_require_middle_name = "0"
		}
		if pf_require_mobile_phone == "" {
			pf_require_mobile_phone = "0"
		}
		if pf_require_email_address == "" {
			pf_require_email_address = "0"
		}
		if pf_require_delivery == "" {
			pf_require_delivery = "0"
		}
		if pf_require_comment == "" {
			pf_require_comment = "0"
		}

		(*wrap.Config).Shop.Price.Format = pfi_price_fomat
		(*wrap.Config).Shop.Price.Round = pfi_price_round

		(*wrap.Config).Shop.Orders.RequiredFields.LastName = utils.StrToInt(pf_require_last_name)
		(*wrap.Config).Shop.Orders.RequiredFields.FirstName = utils.StrToInt(pf_require_first_name)
		(*wrap.Config).Shop.Orders.RequiredFields.MiddleName = utils.StrToInt(pf_require_middle_name)
		(*wrap.Config).Shop.Orders.RequiredFields.MobilePhone = utils.StrToInt(pf_require_mobile_phone)
		(*wrap.Config).Shop.Orders.RequiredFields.EmailAddress = utils.StrToInt(pf_require_email_address)
		(*wrap.Config).Shop.Orders.RequiredFields.Delivery = utils.StrToInt(pf_require_delivery)
		(*wrap.Config).Shop.Orders.RequiredFields.Comment = utils.StrToInt(pf_require_comment)

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
