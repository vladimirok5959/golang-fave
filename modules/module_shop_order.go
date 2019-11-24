package modules

import (
	"strings"

	"golang-fave/engine/basket"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopOrder() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-order",
		WantAdmin: false,
	}, func(wrap *wrapper.Wrapper) {
		if wrap.ShopBasket.ProductsCount(&basket.SBParam{
			R:         wrap.R,
			DB:        wrap.DB,
			Host:      wrap.CurrHost,
			Config:    wrap.Config,
			SessionId: wrap.GetSessionId(),
		}) <= 0 {
			wrap.Write(`{"error": true, "variable": "ShopOrderErrorBasketEmpty"}`)
			return
		}

		pf_client_last_name := wrap.R.FormValue("client_last_name")
		pf_client_first_name := wrap.R.FormValue("client_first_name")
		pf_client_second_name := wrap.R.FormValue("client_second_name")
		pf_client_phone := wrap.R.FormValue("client_phone")
		pf_client_email := wrap.R.FormValue("client_email")
		pf_client_delivery_comment := wrap.R.FormValue("client_delivery_comment")
		pf_client_order_comment := wrap.R.FormValue("client_order_comment")

		if (*wrap.Config).Shop.Orders.RequiredFields.LastName != 0 {
			if strings.TrimSpace(pf_client_last_name) == "" {
				wrap.Write(`{"error": true, "field": "client_last_name", "variable": "ShopOrderEmptyLastName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.FirstName != 0 {
			if strings.TrimSpace(pf_client_first_name) == "" {
				wrap.Write(`{"error": true, "field": "client_first_name", "variable": "ShopOrderEmptyFirstName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.SecondName != 0 {
			if strings.TrimSpace(pf_client_second_name) == "" {
				wrap.Write(`{"error": true, "field": "client_second_name", "variable": "ShopOrderEmptySecondName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.MobilePhone != 0 {
			if strings.TrimSpace(pf_client_phone) == "" || !utils.IsValidMobile(pf_client_phone) {
				wrap.Write(`{"error": true, "field": "client_phone", "variable": "ShopOrderEmptyMobilePhone"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.EmailAddress != 0 {
			if strings.TrimSpace(pf_client_email) == "" || !utils.IsValidEmail(pf_client_email) {
				wrap.Write(`{"error": true, "field": "client_email", "variable": "ShopOrderEmptyEmailAddress"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.Delivery != 0 {
			if strings.TrimSpace(pf_client_delivery_comment) == "" {
				wrap.Write(`{"error": true, "field": "client_delivery_comment", "variable": "ShopOrderEmptyDelivery"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.Comment != 0 {
			if strings.TrimSpace(pf_client_order_comment) == "" {
				wrap.Write(`{"error": true, "field": "client_order_comment", "variable": "ShopOrderEmptyComment"}`)
				return
			}
		}

		// Clear user basket
		wrap.ShopBasket.ClearBasket(&basket.SBParam{
			R:         wrap.R,
			DB:        wrap.DB,
			Host:      wrap.CurrHost,
			Config:    wrap.Config,
			SessionId: wrap.GetSessionId(),
		})

		wrap.Write(`{"error": false, "field": "", "variable": "ShopOrderSuccess"}`)
		return

		// if !utils.IsNumeric(pf_id) {
		// 	wrap.MsgError(`Inner system error`)
		// 	return
		// }

		// if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
		// 	if _, err := tx.Exec(`
		// 		UPDATE shop_products SET
		// 			parent_id = NULL,
		// 			active = 0
		// 		WHERE
		// 			id = ?
		// 		;`,
		// 		utils.StrToInt(pf_id),
		// 	); err != nil {
		// 		return err
		// 	}
		// 	return nil
		// }); err != nil {
		// 	wrap.MsgError(err.Error())
		// 	return
		// }

		// wrap.RecreateProductXmlFile()

		// wrap.ResetCacheBlocks()

		// // Reload current page
		// wrap.Write(`window.location.reload(false);`)
	})
}
