package migrate

import (
	"context"
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/engine/assets/template"
	"golang-fave/engine/sqlw"
)

func Migrate_000000022(ctx context.Context, db *sqlw.DB, host string) error {
	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/maintenance.html", ThemeFiles.AllData["maintenance.html"], 0664); err != nil {
		return err
	}

	return nil
}
