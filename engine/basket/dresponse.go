package basket

import (
	"encoding/json"
)

type dResponse struct {
	IsError bool   `json:"error"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
}

func (this *dResponse) String() string {
	json, err := json.Marshal(this)
	if err != nil {
		return `{"msg":"basket_engine_error","message":"` + err.Error() + `"}`
	}
	return string(json)
}
