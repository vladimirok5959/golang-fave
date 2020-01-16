package modules

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_SettingsShop() *Action {
	return this.newAction(AInfo{
		Mount:     "settings-shop",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_price_fomat := utils.Trim(wrap.R.FormValue("price-fomat"))
		pf_price_round := utils.Trim(wrap.R.FormValue("price-round"))

		pf_require_last_name := utils.Trim(wrap.R.FormValue("require-last-name"))
		pf_require_first_name := utils.Trim(wrap.R.FormValue("require-first-name"))
		pf_require_middle_name := utils.Trim(wrap.R.FormValue("require-middle-name"))
		pf_require_mobile_phone := utils.Trim(wrap.R.FormValue("require-mobile-phone"))
		pf_require_email_address := utils.Trim(wrap.R.FormValue("require-email-address"))
		pf_require_delivery := utils.Trim(wrap.R.FormValue("require-delivery"))
		pf_require_comment := utils.Trim(wrap.R.FormValue("require-comment"))

		pf_new_order_notify_email := utils.Trim(wrap.R.FormValue("new-order-notify-email"))

		pf_new_order_email_theme_cp := utils.Trim(wrap.R.FormValue("new-order-email-theme-cp"))
		pf_new_order_email_theme_user := utils.Trim(wrap.R.FormValue("new-order-email-theme-user"))

		pf_accept_orders := utils.Trim(wrap.R.FormValue("accept-orders"))

		pf_custom_field_1_enabled := utils.Trim(wrap.R.FormValue("custom-field-1-enabled"))
		pf_custom_field_1_caption := utils.Trim(wrap.R.FormValue("custom-field-1-caption"))

		pf_custom_field_2_enabled := utils.Trim(wrap.R.FormValue("custom-field-2-enabled"))
		pf_custom_field_2_caption := utils.Trim(wrap.R.FormValue("custom-field-2-caption"))

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

		if pf_accept_orders == "" {
			pf_accept_orders = "0"
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

		if pf_new_order_notify_email != "" {
			if utils.IsValidEmail(pf_new_order_notify_email) {
				(*wrap.Config).Shop.Orders.NotifyEmail = pf_new_order_notify_email
			}
		} else {
			(*wrap.Config).Shop.Orders.NotifyEmail = ""
		}

		(*wrap.Config).Shop.Orders.NewOrderEmailThemeCp = pf_new_order_email_theme_cp
		(*wrap.Config).Shop.Orders.NewOrderEmailThemeUser = pf_new_order_email_theme_user

		(*wrap.Config).Shop.Orders.Enabled = utils.StrToInt(pf_accept_orders)

		(*wrap.Config).Shop.CustomFields.Field1.Enabled = utils.StrToInt(pf_custom_field_1_enabled)
		(*wrap.Config).Shop.CustomFields.Field1.Caption = pf_custom_field_1_caption

		(*wrap.Config).Shop.CustomFields.Field2.Enabled = utils.StrToInt(pf_custom_field_2_enabled)
		(*wrap.Config).Shop.CustomFields.Field2.Caption = pf_custom_field_2_caption

		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
