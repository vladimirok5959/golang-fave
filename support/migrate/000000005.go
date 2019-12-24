package migrate

import (
	"context"
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/engine/assets/template"
	"golang-fave/engine/sqlw"
)

func Migrate_000000005(ctx context.Context, db *sqlw.DB, host string) error {
	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/shop-category.html", ThemeFiles.AllData["shop-category.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/shop-product.html", ThemeFiles.AllData["shop-product.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/shop.html", ThemeFiles.AllData["shop.html"], 0664); err != nil {
		return err
	}

	return nil
}
