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
		if (*wrap.Config).Shop.Orders.Enabled <= 0 {
			wrap.Write(`{"error": true, "variable": "ShopOrderErrorDisabled"}`)
			return
		}

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
		pf_client_middle_name := wrap.R.FormValue("client_middle_name")
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
		if (*wrap.Config).Shop.Orders.RequiredFields.MiddleName != 0 {
			if strings.TrimSpace(pf_client_middle_name) == "" {
				wrap.Write(`{"error": true, "field": "client_middle_name", "variable": "ShopOrderEmptyMiddleName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.MobilePhone != 0 {
			if strings.TrimSpace(pf_client_phone) == "" {
				wrap.Write(`{"error": true, "field": "client_phone", "variable": "ShopOrderEmptyMobilePhone"}`)
				return
			}
			if !utils.IsValidMobile(pf_client_phone) {
				wrap.Write(`{"error": true, "field": "client_phone", "variable": "ShopOrderNotCorrectMobilePhone"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.EmailAddress != 0 {
			if strings.TrimSpace(pf_client_email) == "" {
				wrap.Write(`{"error": true, "field": "client_email", "variable": "ShopOrderEmptyEmailAddress"}`)
				return
			}
			if !utils.IsValidEmail(pf_client_email) {
				wrap.Write(`{"error": true, "field": "client_email", "variable": "ShopOrderNotCorrectEmailAddress"}`)
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

		bdata := wrap.ShopBasket.GetAll(&basket.SBParam{
			R:         wrap.R,
			DB:        wrap.DB,
			Host:      wrap.CurrHost,
			Config:    wrap.Config,
			SessionId: wrap.GetSessionId(),
		})

		var lastID int64 = 0
		if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Insert row
			res, err := tx.Exec(
				`INSERT INTO shop_orders SET
					create_datetime = ?,
					update_datetime = ?,
					currency_id = ?,
					currency_name = ?,
					currency_coefficient = ?,
					currency_code = ?,
					currency_symbol = ?,
					client_last_name = ?,
					client_first_name = ?,
					client_middle_name = ?,
					client_phone = ?,
					client_email = ?,
					client_delivery_comment = ?,
					client_order_comment = ?,
					status = ?
				;`,
				utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
				utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
				bdata.Currency.Id,
				bdata.Currency.Name,
				bdata.Currency.Coefficient,
				bdata.Currency.Code,
				bdata.Currency.Symbol,
				pf_client_last_name,
				pf_client_first_name,
				pf_client_middle_name,
				pf_client_phone,
				pf_client_email,
				pf_client_delivery_comment,
				pf_client_order_comment,
				0,
			)
			if err != nil {
				return err
			}

			// Get inserted order id
			lastID, err = res.LastInsertId()
			if err != nil {
				return err
			}

			// Insert order products
			for _, product := range *(*bdata).Products {
				if _, err = tx.Exec(
					`INSERT INTO shop_order_products (id, order_id, product_id, price, quantity) VALUES (NULL, ?, ?, ?, ?);`,
					lastID, product.A_product_id, product.A_price, product.A_quantity,
				); err != nil {
					return err
				}
			}

			// Send notify email
			if (*wrap.Config).Shop.Orders.NotifyEmail != "" {
				if err := wrap.SendEmail(
					(*wrap.Config).Shop.Orders.NotifyEmail,
					"❤️ New Order ("+wrap.Host+":"+wrap.Port+")",
					"You have new order in shop on host: <a href=\"http://"+wrap.Host+":"+wrap.Port+"/\">http://"+wrap.Host+":"+wrap.Port+"/</a>",
				); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			wrap.Write(`{"error": true, "variable": "ShopOrderErrorSomethingWrong"}`)
			return
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
	})
}
