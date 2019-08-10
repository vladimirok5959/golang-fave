package migrate

import (
	"fmt"

	"golang-fave/engine/sqlw"
	"golang-fave/utils"
)

func Run(db *sqlw.DB, version int, host string) error {
	var last string
	for i, fn := range Migrations {
		if utils.StrToInt(i) > 1 {
			if version < utils.StrToInt(i) {
				last = i
				if fn != nil {
					fn(db, host)
					fmt.Printf("Migrated %s: %s\n", host, i)
				}
			}
		}
	}

	if last != "" {
		if _, err := db.Exec(`UPDATE settings SET value = ? WHERE name = 'database_version';`, last); err != nil {
			return err
		}
	}

	return nil
}
