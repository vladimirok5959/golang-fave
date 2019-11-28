package utils

import (
	"database/sql"
)

type MySql_shop_product struct {
	A_id        int
	A_parent    sql.NullInt64
	A_user      int
	A_currency  int
	A_price     float64
	A_price_old float64
	A_gname     string
	A_name      string
	A_alias     string
	A_vendor    string
	A_quantity  int
	A_category  int
	A_briefly   string
	A_content   string
	A_datetime  int
	A_active    int
}

func (this *MySql_shop_product) A_parent_id() int {
	return int(this.A_parent.Int64)
}
