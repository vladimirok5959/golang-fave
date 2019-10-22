package basket

type product struct {
	currency *currency
	price    float64

	Id       int    `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Link     string `json:"link"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
	Sum      string `json:"sum"`
}
