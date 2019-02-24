package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"html"
	"os"
	"strconv"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterModule_Index() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "index",
		Name:   "Pages",
		Order:  0,
		Icon:   assets.SysSvgIconPage,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of Pages", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add New Page", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify Page", Show: false},
		},
	}, func(wrap *wrapper.Wrapper) {
		// Front-end
		wrap.RenderFrontEnd("index", consts.TmplDataModIndex{
			MetaTitle:       "Meta Title",
			MetaKeywords:    "Meta Keywords",
			MetaDescription: "Meta Description",

			MainMenuItems: []consts.TmplDataMainMenuItem{
				{Name: "Home", Link: "/", Active: true},
				{Name: "Item 1", Link: "/#1", Active: false},
				{Name: "Item 2", Link: "/#2", Active: false},
				{Name: "Item 3", Link: "/#3", Active: false},
			},
		})
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of Pages"},
			})
			content += builder.DataTable(wrap, "pages", "id", "DESC", []builder.DataTableRow{
				{
					DBField: "id",
				},
				{
					DBField:     "name",
					NameInTable: "Page / Alias",
					CallBack: func(values *[]string) string {
						name := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[1]) + `</a>`
						alias := html.EscapeString((*values)[2])
						return `<div>` + name + `</div><div><small>` + alias + `</small></div>`
					},
				},
				{
					DBField: "alias",
				},
				{
					DBField:     "datetime",
					NameInTable: "Date / Time",
				},
				{
					DBField:     "active",
					NameInTable: "Active",
				},
			}, func(values *[]string) string {
				return `<a class="ico" href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` +
					assets.SysSvgIconEdit + `</a>` +
					`<a class="ico" href="#">` + assets.SysSvgIconRemove + `</a>`
			}, "/cp/"+wrap.CurrModule+"/")
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "Add New Page"},
			})
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}

func (this *Modules) RegisterAction_IndexMysqlSetup() *Action {
	return this.newAction(AInfo{
		WantDB: false,
		Mount:  "index-mysql-setup",
	}, func(wrap *wrapper.Wrapper) {
		pf_host := wrap.R.FormValue("host")
		pf_port := wrap.R.FormValue("port")
		pf_name := wrap.R.FormValue("name")
		pf_user := wrap.R.FormValue("user")
		pf_password := wrap.R.FormValue("password")

		if pf_host == "" {
			wrap.MsgError(`Please specify host for MySQL connection`)
			return
		}

		if pf_port == "" {
			wrap.MsgError(`Please specify host port for MySQL connection`)
			return
		}

		if _, err := strconv.Atoi(pf_port); err != nil {
			wrap.MsgError(`MySQL host port must be integer number`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify MySQL database name`)
			return
		}

		if pf_user == "" {
			wrap.MsgError(`Please specify MySQL user`)
			return
		}

		// Try connect to mysql
		db, err := sql.Open("mysql", pf_user+":"+pf_password+"@tcp("+pf_host+":"+pf_port+")/"+pf_name)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Try to install all tables
		_, err = db.Query(fmt.Sprintf(
			`CREATE TABLE %s.users (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				first_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'User first name',
				last_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'User last name',
				email VARCHAR(64) NOT NULL COMMENT 'User email',
				password VARCHAR(32) NOT NULL COMMENT 'User password (MD5)',
				admin int(1) NOT NULL COMMENT 'Is admin user or not',
				active int(1) NOT NULL COMMENT 'Is active user or not',
				PRIMARY KEY (id)
			) ENGINE = InnoDB;`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`CREATE TABLE %s.pages (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				user int(11) NOT NULL COMMENT 'User id',
				name varchar(255) NOT NULL COMMENT 'Page name',
				alias varchar(255) NOT NULL COMMENT 'Page url part',
				content text NOT NULL COMMENT 'Page content',
				meta_title varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta title',
				meta_keywords varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta keywords',
				meta_description varchar(510) NOT NULL DEFAULT '' COMMENT 'Page meta description',
				datetime datetime NOT NULL COMMENT 'Creation date/time',
				active int(1) NOT NULL COMMENT 'Is active page or not',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Save mysql config file
		err = utils.MySqlConfigWrite(wrap.DConfig+string(os.PathSeparator)+"mysql.json", pf_host, pf_port, pf_name, pf_user, pf_password)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}

func (this *Modules) RegisterAction_IndexFirstUser() *Action {
	return this.newAction(AInfo{
		WantDB: true,
		Mount:  "index-first-user",
	}, func(wrap *wrapper.Wrapper) {
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		if pf_password == "" {
			wrap.MsgError(`Please specify user password`)
			return
		}

		_, err := wrap.DB.Query(
			`INSERT INTO users SET
				first_name = ?,
				last_name = ?,
				email = ?,
				password = MD5(?),
				admin = 1,
				active = 1
			;`,
			pf_first_name,
			pf_last_name,
			pf_email,
			pf_password,
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}

func (this *Modules) RegisterAction_IndexUserSignIn() *Action {
	return this.newAction(AInfo{
		WantDB: true,
		Mount:  "index-user-sign-in",
	}, func(wrap *wrapper.Wrapper) {
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		if pf_password == "" {
			wrap.MsgError(`Please specify user password`)
			return
		}

		if wrap.S.GetInt("UserId", 0) > 0 {
			wrap.MsgError(`You already logined`)
			return
		}

		var user_id int
		err := wrap.DB.QueryRow(
			`SELECT
				id
			FROM
				users
			WHERE
				email = ? and
				password = MD5(?) and
				admin = 1 and
				active = 1
			LIMIT 1;`,
			pf_email,
			pf_password,
		).Scan(
			&user_id,
		)

		if err != nil && err != sql.ErrNoRows {
			wrap.MsgError(err.Error())
			return
		}

		if err == sql.ErrNoRows {
			wrap.MsgError(`Incorrect email or password`)
			return
		}

		// Save to current session
		wrap.S.SetInt("UserId", user_id)

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}

func (this *Modules) RegisterAction_CpUserLogout() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "singout",
		WantUser: true,
	}, func(wrap *wrapper.Wrapper) {
		// Reset session var
		wrap.S.SetInt("UserId", 0)

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}

func (this *Modules) RegisterAction_CpUserSettings() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "user-settings",
		WantUser: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_first_name := wrap.R.FormValue("first_name")
		pf_last_name := wrap.R.FormValue("last_name")
		pf_email := wrap.R.FormValue("email")
		pf_password := wrap.R.FormValue("password")

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		if pf_password != "" {
			// Update with password if set
			_, err := wrap.DB.Query(
				`UPDATE users SET
					first_name = ?,
					last_name = ?,
					email = ?,
					password = MD5(?)
				WHERE
					id = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				pf_password,
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		} else {
			// Update without password if not set
			_, err := wrap.DB.Query(
				`UPDATE users SET
					first_name = ?,
					last_name = ?,
					email = ?
				WHERE
					id = ?
				;`,
				pf_first_name,
				pf_last_name,
				pf_email,
				wrap.User.A_id,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
