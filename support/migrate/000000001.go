package migrate

import (
	"golang-fave/engine/sqlw"
)

var Migrations = map[string]func(*sqlw.DB) error{
	"000000000": nil,
	"000000001": nil,
}
