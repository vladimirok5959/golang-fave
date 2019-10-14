package basket

type currency struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Coefficient float64 `json:"coefficient"`
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
}
