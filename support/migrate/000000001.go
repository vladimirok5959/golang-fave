package migrate

import (
	"golang-fave/engine/sqlw"
)

var Migrations = map[string]func(*sqlw.DB, string) error{
	"000000000": nil,
	"000000001": nil,
	"000000002": Migrate_000000002,
	"000000003": Migrate_000000003,
	"000000004": Migrate_000000004,
}
