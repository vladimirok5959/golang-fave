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

	RenderName     string
	RenderLink     string
	RenderPrice    string
	RenderQuantity int
	RenderSum      string
}

type MySql_basket struct {
	Products   *[]MySql_basket_product
	Currency   *MySql_basket_currency
	TotalSum   float64
	TotalCount int

	RenderTotalSum string
}
