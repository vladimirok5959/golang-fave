package basket

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Sum      float64 `json:"sum"`
}
