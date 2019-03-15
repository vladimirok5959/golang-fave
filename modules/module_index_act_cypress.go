package modules

import (
	"os"

	"golang-fave/consts"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_IndexCypressReset() *Action {
	return this.newAction(AInfo{
		WantDB: true,
		Mount:  "index-cypress-reset",
	}, func(wrap *wrapper.Wrapper) {
		if !consts.ParamDebug {
			wrap.Write(`Access denied`)
			return
		}

		_, err := wrap.DB.Query(
			`DROP TABLE
				blog_cats,
				blog_cat_post_rel,
				blog_posts,
				pages,
				users
			;`,
		)
		if err != nil {
			wrap.Write(err.Error())
			return
		}

		err = os.Remove(wrap.DConfig + string(os.PathSeparator) + ".installed")
		if err != nil {
			wrap.Write(err.Error())
			return
		}

		err = os.Remove(wrap.DConfig + string(os.PathSeparator) + "mysql.json")
		if err != nil {
			wrap.Write(err.Error())
			return
		}

		wrap.Write(`OK`)
	})
}
