package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"

	"golang-fave/assets"
	"golang-fave/consts"
	"golang-fave/engine/builder"
	"golang-fave/engine/fetdata"
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
			{Mount: "default", Name: "List of pages", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new page", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify page", Show: false},
		},
	}, func(wrap *wrapper.Wrapper) {
		// Front-end
		row := &utils.MySql_page{}
		err := wrap.DB.QueryRow(`
			SELECT
				id,
				user,
				name,
				alias,
				content,
				meta_title,
				meta_keywords,
				meta_description,
				UNIX_TIMESTAMP(datetime) as datetime,
				active
			FROM
				pages
			WHERE
				active = 1 and
				alias = ?
			LIMIT 1;`,
			wrap.R.URL.Path,
		).Scan(
			&row.A_id,
			&row.A_user,
			&row.A_name,
			&row.A_alias,
			&row.A_content,
			&row.A_meta_title,
			&row.A_meta_keywords,
			&row.A_meta_description,
			&row.A_datetime,
			&row.A_active,
		)
		if err != nil && err != sql.ErrNoRows {
			// System error 500
			utils.SystemErrorPageEngine(wrap.W, err)
			return
		} else if err == sql.ErrNoRows {
			// User error 404 page
			wrap.W.WriteHeader(http.StatusNotFound)
			wrap.RenderFrontEnd("404", fetdata.New(wrap, nil, true))
			return
		}

		// Replace title with page name
		if row.A_meta_title == "" {
			row.A_meta_title = row.A_name
		}

		// Which template
		tmpl_name := "index"
		if wrap.R.URL.Path != "/" {
			tmpl_name = "page"
		}

		// Render template
		wrap.RenderFrontEnd(tmpl_name, fetdata.New(wrap, row, false))
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of pages"},
			})
			content += builder.DataTable(wrap, "pages", "id", "DESC", &[]builder.DataTableRow{
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
					DBExp:       "UNIX_TIMESTAMP(`datetime`)",
					NameInTable: "Date / Time",
					Classes:     "d-none d-md-table-cell",
					CallBack: func(values *[]string) string {
						t := int64(utils.StrToInt((*values)[3]))
						return `<div>` + utils.UnixTimestampToFormat(t, "02.01.2006") + `</div>` +
							`<div><small>` + utils.UnixTimestampToFormat(t, "15:04:05") + `</small></div>`
					},
				},
				{
					DBField:     "active",
					NameInTable: "Active",
					Classes:     "d-none d-sm-table-cell",
					CallBack: func(values *[]string) string {
						return builder.CheckBox(utils.StrToInt((*values)[4]))
					},
				},
			}, func(values *[]string) string {
				return builder.DataTableAction(&[]builder.DataTableActionRow{
					{
						Icon:   assets.SysSvgIconView,
						Href:   (*values)[2],
						Hint:   "View",
						Target: "_blank",
					},
					{
						Icon: assets.SysSvgIconEdit,
						Href: "/cp/" + wrap.CurrModule + "/modify/" + (*values)[0] + "/",
						Hint: "Edit",
					},
					{
						Icon: assets.SysSvgIconRemove,
						Href: "javascript:fave.ActionDataTableDelete(this,'index-delete','" +
							(*values)[0] + "','Are you sure want to delete page?');",
						Hint:    "Delete",
						Classes: "delete",
					},
				})
			}, "/cp/"+wrap.CurrModule+"/")
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new page"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify page"},
				})
			}

			data := utils.MySql_page{
				A_id:               0,
				A_user:             0,
				A_name:             "",
				A_alias:            "",
				A_content:          "",
				A_meta_title:       "",
				A_meta_keywords:    "",
				A_meta_description: "",
				A_datetime:         0,
				A_active:           0,
			}

			if wrap.CurrSubModule == "modify" {
				if len(wrap.UrlArgs) != 3 {
					return "", "", ""
				}
				if !utils.IsNumeric(wrap.UrlArgs[2]) {
					return "", "", ""
				}
				err := wrap.DB.QueryRow(`
					SELECT
						id,
						user,
						name,
						alias,
						content,
						meta_title,
						meta_keywords,
						meta_description,
						active
					FROM
						pages
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_name,
					&data.A_alias,
					&data.A_content,
					&data.A_meta_title,
					&data.A_meta_keywords,
					&data.A_meta_description,
					&data.A_active,
				)
				if err != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "index-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page name",
					Name:    "name",
					Value:   data.A_name,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: /about-us/ or /about-us.html",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Page content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "autosize",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta title",
					Name:    "meta_title",
					Value:   data.A_meta_title,
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta keywords",
					Name:    "meta_keywords",
					Value:   data.A_meta_keywords,
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Meta description",
					Name:    "meta_description",
					Value:   data.A_meta_description,
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Active",
					Name:    "active",
					Value:   utils.IntToStr(data.A_active),
				},
				{
					Kind: builder.DFKMessage,
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  btn_caption,
					Target: "add-edit-button",
				},
			})

			if wrap.CurrSubModule == "add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}

func (this *Modules) RegisterAction_IndexModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "index-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_content := wrap.R.FormValue("content")
		pf_meta_title := wrap.R.FormValue("meta_title")
		pf_meta_keywords := wrap.R.FormValue("meta_keywords")
		pf_meta_description := wrap.R.FormValue("meta_description")
		pf_active := wrap.R.FormValue("active")

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify page name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateAlias(pf_name)
		}

		if !utils.IsValidAlias(pf_alias) {
			wrap.MsgError(`Please specify correct page alias`)
			return
		}

		if pf_id == "0" {
			// Add new page
			_, err := wrap.DB.Query(
				`INSERT INTO pages SET
					user = ?,
					name = ?,
					alias = ?,
					content = ?,
					meta_title = ?,
					meta_keywords = ?,
					meta_description = ?,
					datetime = ?,
					active = ?
				;`,
				wrap.User.A_id,
				pf_name,
				pf_alias,
				pf_content,
				pf_meta_title,
				pf_meta_keywords,
				pf_meta_description,
				utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
				pf_active,
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/';`)
		} else {
			// Update page
			_, err := wrap.DB.Query(
				`UPDATE pages SET
					name = ?,
					alias = ?,
					content = ?,
					meta_title = ?,
					meta_keywords = ?,
					meta_description = ?,
					active = ?
				WHERE
					id = ?
				;`,
				pf_name,
				pf_alias,
				pf_content,
				pf_meta_title,
				pf_meta_keywords,
				pf_meta_description,
				pf_active,
				utils.StrToInt(pf_id),
			)
			if err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.Write(`window.location='/cp/index/modify/` + pf_id + `/';`)
		}
	})
}

func (this *Modules) RegisterAction_IndexDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "index-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Delete page
		_, err := wrap.DB.Query(
			`DELETE FROM pages WHERE id = ?;`,
			utils.StrToInt(pf_id),
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
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

		// Security, check if still need to run this action
		if wrap.ConfMysqlExists {
			wrap.MsgError(`MySQL is already configured`)
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
			`ALTER TABLE %s.users ADD UNIQUE KEY email (email);`,
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
		_, err = db.Query(fmt.Sprintf(
			`ALTER TABLE %s.pages ADD UNIQUE KEY alias (alias);`,
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

		// Security, check if still need to run this action
		var count int
		err := wrap.DB.QueryRow(`
			SELECT
				COUNT(*)
			FROM
				users
			;`,
		).Scan(
			&count,
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		if count > 0 {
			wrap.MsgError(`MySQL is already configured`)
			return
		}

		_, err = wrap.DB.Query(
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

		// Add home page
		_, err = wrap.DB.Query(
			`INSERT INTO pages SET
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			"Home",
			"/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Add another page
		_, err = wrap.DB.Query(
			`INSERT INTO pages SET
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			"Another",
			"/another/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		)
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Add about page
		_, err = wrap.DB.Query(
			`INSERT INTO pages SET
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			"About",
			"/about/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
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

func (this *Modules) RegisterAction_IndexUserLogout() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "index-user-logout",
		WantUser: true,
	}, func(wrap *wrapper.Wrapper) {
		// Reset session var
		wrap.S.SetInt("UserId", 0)

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}

func (this *Modules) RegisterAction_IndexUserUpdateProfile() *Action {
	return this.newAction(AInfo{
		WantDB:   true,
		Mount:    "index-user-update-profile",
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
