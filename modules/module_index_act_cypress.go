package modules

import (
	"os"

	"golang-fave/consts"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_IndexCypressReset() *Action {
	return this.newAction(AInfo{
		WantDB: false,
		Mount:  "index-cypress-reset",
	}, func(wrap *wrapper.Wrapper) {
		if !consts.ParamDebug {
			wrap.Write(`Access denied`)
			return
		}

		db, err := sqlw.Open("mysql", "root:root@tcp(localhost:3306)/fave")
		if err != nil {
			wrap.Write(err.Error())
			return
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			wrap.Write(err.Error())
			return
		}

		os.Remove(wrap.DConfig + string(os.PathSeparator) + ".installed")
		os.Remove(wrap.DConfig + string(os.PathSeparator) + "mysql.json")

		_, _ = db.Exec(
			`DROP TABLE
				blog_cats,
				blog_cat_post_rel,
				blog_posts,
				pages,
				users
			;`,
		)

		wrap.Write(`OK`)
	})
}
