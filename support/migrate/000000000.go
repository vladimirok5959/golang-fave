package migrate

import (
	"context"
	"fmt"

	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"
)

func Run(ctx context.Context, db *sqlw.DB, version int, host string) error {
	var last string
	for i, fn := range Migrations {
		if utils.StrToInt(i) > 1 {
			if version < utils.StrToInt(i) {
				last = i
				if fn != nil {
					fn(ctx, db, host)
					fmt.Printf("Migrated %s: %s\n", host, i)
				}
			}
		}
	}

	if last != "" {
		if _, err := db.Exec(ctx, `UPDATE fave_settings SET value = ? WHERE name = 'database_version';`, last); err != nil {
			return err
		}
	}

	return nil
}
