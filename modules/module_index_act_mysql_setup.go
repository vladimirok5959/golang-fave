package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"os"
	"strconv"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

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
			wrap.MsgError(`CMS is already configured`)
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
		_, err = wrap.DB.Query(
			`INSERT INTO %s.pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			pf_name,
			1,
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
		_, err = wrap.DB.Query(
			`INSERT INTO %s.pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			pf_name,
			2,
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
		_, err = wrap.DB.Query(
			`INSERT INTO %s.pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			pf_name,
			3,
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
		_, err = db.Query(fmt.Sprintf(
			`ALTER TABLE %s.pages ADD UNIQUE KEY alias (alias);`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`CREATE TABLE %s.blog_posts (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				user int(11) NOT NULL COMMENT 'User id',
				name varchar(255) NOT NULL COMMENT 'Post name',
				alias varchar(255) NOT NULL COMMENT 'Post alias',
				content text NOT NULL COMMENT 'Post content',
				datetime datetime NOT NULL COMMENT 'Creation date/time',
				active int(1) NOT NULL COMMENT 'Is active post or not',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`ALTER TABLE %s.blog_posts ADD UNIQUE KEY alias (alias);`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`CREATE TABLE %s.blog_cats (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				user int(11) NOT NULL COMMENT 'User id',
				name varchar(255) NOT NULL COMMENT 'Category name',
				alias varchar(255) NOT NULL COMMENT 'Category alias',
				lft int(11) NOT NULL COMMENT 'For nested set model',
				rgt int(11) NOT NULL COMMENT 'For nested set model',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`ALTER TABLE %s.blog_cats ADD UNIQUE KEY alias (alias);`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`INSERT INTO %s.blog_cats (id, user, name, alias, lft, rgt) VALUES (1, 0, 'ROOT', 'ROOT', 1, 2);`,
			pf_name))
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}
		_, err = db.Query(fmt.Sprintf(
			`ALTER TABLE %s.blog_cats ADD KEY lft (lft), ADD KEY rgt (rgt);`,
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
