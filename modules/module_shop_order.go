package modules

import (
	"golang-fave/engine/wrapper"
	// "golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopOrder() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-order",
		WantAdmin: false,
	}, func(wrap *wrapper.Wrapper) {
		// pf_id := wrap.R.FormValue("id")

		// pf_client_last_name := wrap.R.FormValue("client_last_name")
		// pf_client_first_name := wrap.R.FormValue("client_first_name")
		// pf_client_second_name := wrap.R.FormValue("client_second_name")
		// pf_client_phone := wrap.R.FormValue("client_phone")
		// pf_client_email := wrap.R.FormValue("client_email")
		// pf_client_delivery_comment := wrap.R.FormValue("client_delivery_comment")
		// pf_client_order_comment := wrap.R.FormValue("client_order_comment")

		wrap.MsgError(`OK!`)
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
