package support

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"golang-fave/engine/sqlw"
	"golang-fave/support/migrate"
	"golang-fave/utils"
)

type Support struct {
}

func New() *Support {
	sup := Support{}
	return &sup
}

func (this *Support) isSettingsTableDoesntExist(err error) bool {
	error_msg := strings.ToLower(err.Error())
	if match, _ := regexp.MatchString(`^error 1146`, error_msg); match {
		if match, _ := regexp.MatchString(`'[^\.]+\.fave_settings'`, error_msg); match {
			if match, _ := regexp.MatchString(`doesn't exist$`, error_msg); match {
				return true
			}
		}
	}
	return false
}

func (this *Support) Migration(ctx context.Context, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if utils.IsDir(dir + string(os.PathSeparator) + file.Name()) {
			if err := this.Migrate(ctx, dir+string(os.PathSeparator)+file.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Support) Migrate(ctx context.Context, host string) error {
	mysql_config_file := host + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "mysql.json"
	if utils.IsMySqlConfigExists(mysql_config_file) {
		mc, err := utils.MySqlConfigRead(mysql_config_file)
		if err != nil {
			return err
		}
		db, err := sqlw.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
		if err != nil {
			return err
		}
		if err := db.Ping(ctx); err != nil {
			return err
		}
		defer db.Close()

		var table string
		if err := db.QueryRow(ctx, `SHOW TABLES LIKE 'settings';`).Scan(&table); err == nil {
			if table == "settings" {
				if _, err := db.Exec(ctx, `RENAME TABLE settings TO fave_settings;`); err != nil {
					return err
				}
			}
		}

		var version string
		if err := db.QueryRow(ctx, `SELECT value FROM fave_settings WHERE name = 'database_version' LIMIT 1;`).Scan(&version); err != nil {
			if this.isSettingsTableDoesntExist(err) {
				if _, err := db.Exec(
					ctx,
					`CREATE TABLE fave_settings (
						name varchar(255) NOT NULL COMMENT 'Setting name',
						value text NOT NULL COMMENT 'Setting value'
					) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
				); err != nil {
					return err
				}
				if _, err := db.Exec(
					ctx,
					`INSERT INTO fave_settings (name, value) VALUES ('database_version', '000000002');`,
				); err != nil {
					return err
				}
				version = "000000002"
				err = nil
			}
			return err
		}
		return this.Process(ctx, db, version, host)
	}
	return nil
}

func (this *Support) Process(ctx context.Context, db *sqlw.DB, version string, host string) error {
	return migrate.Run(ctx, db, utils.StrToInt(version), host)
}
