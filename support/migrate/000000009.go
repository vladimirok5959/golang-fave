package migrate

import (
	"context"
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/engine/assets/template"
	"golang-fave/engine/sqlw"
)

func Migrate_000000009(ctx context.Context, db *sqlw.DB, host string) error {
	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/cached-block-1.html", ThemeFiles.AllData["cached-block-1.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/cached-block-2.html", ThemeFiles.AllData["cached-block-2.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/cached-block-3.html", ThemeFiles.AllData["cached-block-3.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/cached-block-4.html", ThemeFiles.AllData["cached-block-4.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/cached-block-5.html", ThemeFiles.AllData["cached-block-5.html"], 0664); err != nil {
		return err
	}

	return nil
}
