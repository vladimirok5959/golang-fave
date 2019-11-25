package utils

type MySql_basket_currency struct {
	Id          int
	Name        string
	Coefficient float64
	Code        string
	Symbol      string
}

type MySql_basket_product struct {
	A_product_id int
	A_price      float64
	A_quantity   int
}

type MySql_basket_order struct {
	Products   *[]MySql_basket_product
	Currency   *MySql_basket_currency
	TotalSum   float64
	TotalCount int
}

// type MySql_basket_order struct {
// 	A_currency_id int
// 	A_currency_name string
// 	A_currency_coefficient float64
// 	A_currency_code string
// 	A_currency_symbol string
// 	A_client_last_name string
// 	A_client_first_name string
// 	A_client_second_name string
// 	A_client_phone string
// 	A_client_email string
// 	A_client_delivery_comment string
// 	A_client_order_comment string
// 	A_status string
// }
