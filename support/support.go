package support

import (
	"io/ioutil"
	"os"

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

func (this *Support) Migration(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if utils.IsDir(dir + string(os.PathSeparator) + file.Name()) {
			if err := this.Migrate(dir + string(os.PathSeparator) + file.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Support) Migrate(host string) error {
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
		if err := db.Ping(); err != nil {
			return err
		}
		defer db.Close()
		var version string
		if err := db.QueryRow(`SELECT value FROM settings WHERE name = 'database_version' LIMIT 1;`).Scan(&version); err != nil {
			return err
		}
		return this.Process(db, version, host)
	}
	return nil
}

func (this *Support) Process(db *sqlw.DB, version string, host string) error {
	return migrate.Run(db, utils.StrToInt(version), host)
}
