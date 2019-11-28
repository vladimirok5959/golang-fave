package modules

import (
	"strings"

	"golang-fave/consts"
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

		pf_client_last_name := strings.TrimSpace(wrap.R.FormValue("client_last_name"))
		pf_client_first_name := strings.TrimSpace(wrap.R.FormValue("client_first_name"))
		pf_client_middle_name := strings.TrimSpace(wrap.R.FormValue("client_middle_name"))
		pf_client_phone := strings.TrimSpace(wrap.R.FormValue("client_phone"))
		pf_client_email := strings.TrimSpace(wrap.R.FormValue("client_email"))
		pf_client_delivery_comment := strings.TrimSpace(wrap.R.FormValue("client_delivery_comment"))
		pf_client_order_comment := strings.TrimSpace(wrap.R.FormValue("client_order_comment"))

		if (*wrap.Config).Shop.Orders.RequiredFields.LastName != 0 {
			if pf_client_last_name == "" {
				wrap.Write(`{"error": true, "field": "client_last_name", "variable": "ShopOrderEmptyLastName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.FirstName != 0 {
			if pf_client_first_name == "" {
				wrap.Write(`{"error": true, "field": "client_first_name", "variable": "ShopOrderEmptyFirstName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.MiddleName != 0 {
			if pf_client_middle_name == "" {
				wrap.Write(`{"error": true, "field": "client_middle_name", "variable": "ShopOrderEmptyMiddleName"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.MobilePhone != 0 {
			if pf_client_phone == "" {
				wrap.Write(`{"error": true, "field": "client_phone", "variable": "ShopOrderEmptyMobilePhone"}`)
				return
			}
			if !utils.IsValidMobile(pf_client_phone) {
				wrap.Write(`{"error": true, "field": "client_phone", "variable": "ShopOrderNotCorrectMobilePhone"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.EmailAddress != 0 {
			if pf_client_email == "" {
				wrap.Write(`{"error": true, "field": "client_email", "variable": "ShopOrderEmptyEmailAddress"}`)
				return
			}
			if !utils.IsValidEmail(pf_client_email) {
				wrap.Write(`{"error": true, "field": "client_email", "variable": "ShopOrderNotCorrectEmailAddress"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.Delivery != 0 {
			if pf_client_delivery_comment == "" {
				wrap.Write(`{"error": true, "field": "client_delivery_comment", "variable": "ShopOrderEmptyDelivery"}`)
				return
			}
		}
		if (*wrap.Config).Shop.Orders.RequiredFields.Comment != 0 {
			if pf_client_order_comment == "" {
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

			// Send notify email to owner
			if (*wrap.Config).Shop.Orders.NotifyEmail != "" {
				if err := wrap.SendEmailTemplated(
					(*wrap.Config).Shop.Orders.NotifyEmail,
					(*wrap.Config).Shop.Orders.NewOrderEmailThemeCp+" #"+utils.Int64ToStr(lastID)+" ("+wrap.Host+":"+wrap.Port+")",
					"email-new-order-admin",
					consts.TmplEmailOrder{
						Basket: bdata,
						Client: consts.TmplOrderClient{
							LastName:        pf_client_last_name,
							FirstName:       pf_client_first_name,
							MiddleName:      pf_client_middle_name,
							Phone:           pf_client_phone,
							Email:           pf_client_email,
							DeliveryComment: pf_client_delivery_comment,
							OrderComment:    pf_client_order_comment,
						},
						Else: consts.TmplOrderElse{
							OrderId:     lastID,
							Subject:     (*wrap.Config).Shop.Orders.NewOrderEmailThemeCp + " #" + utils.Int64ToStr(lastID) + " (" + wrap.Host + ":" + wrap.Port + ")",
							CpOrderLink: "http://" + wrap.Host + ":" + wrap.Port + "/cp/shop/orders-modify/" + utils.Int64ToStr(lastID) + "/",
						},
					},
				); err != nil {
					return err
				}
			}

			// Send notify email to client
			if pf_client_email != "" {
				if err := wrap.SendEmailTemplated(
					pf_client_email,
					(*wrap.Config).Shop.Orders.NewOrderEmailThemeUser+" #"+utils.Int64ToStr(lastID),
					"email-new-order-user",
					consts.TmplEmailOrder{
						Basket: bdata,
						Client: consts.TmplOrderClient{
							LastName:        pf_client_last_name,
							FirstName:       pf_client_first_name,
							MiddleName:      pf_client_middle_name,
							Phone:           pf_client_phone,
							Email:           pf_client_email,
							DeliveryComment: pf_client_delivery_comment,
							OrderComment:    pf_client_order_comment,
						},
						Else: consts.TmplOrderElse{
							OrderId:     lastID,
							Subject:     (*wrap.Config).Shop.Orders.NewOrderEmailThemeUser + " #" + utils.Int64ToStr(lastID),
							CpOrderLink: "http://" + wrap.Host + ":" + wrap.Port + "/cp/shop/orders-modify/" + utils.Int64ToStr(lastID) + "/",
						},
					},
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
