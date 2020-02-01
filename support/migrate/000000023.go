package migrate

import (
	"context"
	"os"

	"golang-fave/engine/sqlw"
)

func Migrate_000000023(ctx context.Context, db *sqlw.DB, host string) error {
	if err := os.Mkdir(host+string(os.PathSeparator)+"/htdocs/public", 0755); err != nil {
		return err
	}

	return nil
}
