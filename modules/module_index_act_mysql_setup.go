package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

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

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Table: blog_cats
		if _, err = tx.Exec(
			`CREATE TABLE blog_cats (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				user int(11) NOT NULL COMMENT 'User id',
				name varchar(255) NOT NULL COMMENT 'Category name',
				alias varchar(255) NOT NULL COMMENT 'Category alias',
				lft int(11) NOT NULL COMMENT 'For nested set model',
				rgt int(11) NOT NULL COMMENT 'For nested set model',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`INSERT INTO blog_cats (id, user, name, alias, lft, rgt) VALUES (1, 0, 'ROOT', 'ROOT', 1, 2);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO blog_cats (id, user, name, alias, lft, rgt)
				VALUES
			(1, 0, 'ROOT', 'ROOT', 1, 24),
			(2, 1, 'Health and food', 'health-and-food', 2, 15),
			(3, 1, 'News', 'news', 16, 21),
			(4, 1, 'Hobby', 'hobby', 22, 23),
			(5, 1, 'Juices', 'juices', 3, 8),
			(6, 1, 'Nutrition', 'nutrition', 9, 14),
			(7, 1, 'Natural', 'natural', 4, 5),
			(8, 1, 'For kids', 'for-kids', 6, 7),
			(9, 1, 'For all', 'for-all', 10, 11),
			(10, 1, 'For athletes', 'for-athletes', 12, 13),
			(11, 1, 'Computers and technology', 'computers-and-technology', 17, 18),
			(13, 1, 'Film industry', 'film-industry', 19, 20);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cats ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cats ADD KEY lft (lft), ADD KEY rgt (rgt);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: blog_cat_post_rel
		if _, err = tx.Exec(
			`CREATE TABLE blog_cat_post_rel (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				post_id int(11) NOT NULL COMMENT 'Post id',
				category_id int(11) NOT NULL COMMENT 'Category id',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cat_post_rel ADD KEY post_id (post_id), ADD KEY category_id (category_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cat_post_rel ADD UNIQUE KEY post_category (post_id,category_id) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: blog_posts
		if _, err = tx.Exec(
			`CREATE TABLE blog_posts (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				user int(11) NOT NULL COMMENT 'User id',
				name varchar(255) NOT NULL COMMENT 'Post name',
				alias varchar(255) NOT NULL COMMENT 'Post alias',
				content text NOT NULL COMMENT 'Post content',
				datetime datetime NOT NULL COMMENT 'Creation date/time',
				active int(1) NOT NULL COMMENT 'Is active post or not',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_posts ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: pages
		if _, err = tx.Exec(
			`CREATE TABLE pages (
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
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			1,
			"Home",
			"/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			2,
			1,
			"Another",
			"/another/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO pages SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			3,
			1,
			"About",
			"/about/",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p><p>Ante metus dictum at tempor commodo ullamcorper a. Facilisis mauris sit amet massa vitae. Enim neque volutpat ac tincidunt vitae. Tempus quam pellentesque nec nam aliquam sem. Mollis aliquam ut porttitor leo a diam sollicitudin. Nunc pulvinar sapien et ligula ullamcorper. Dignissim suspendisse in est ante in nibh mauris. Eget egestas purus viverra accumsan in. Vitae tempus quam pellentesque nec nam aliquam sem et. Sodales ut etiam sit amet nisl. Aliquet risus feugiat in ante. Rhoncus urna neque viverra justo nec ultrices dui sapien. Sit amet aliquam id diam maecenas ultricies. Sed odio morbi quis commodo odio aenean sed adipiscing diam.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE pages ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: users
		if _, err = tx.Exec(
			`CREATE TABLE users (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				first_name varchar(64) NOT NULL DEFAULT '' COMMENT 'User first name',
				last_name varchar(64) NOT NULL DEFAULT '' COMMENT 'User last name',
				email varchar(64) NOT NULL COMMENT 'User email',
				password varchar(32) NOT NULL COMMENT 'User password (MD5)',
				admin int(1) NOT NULL COMMENT 'Is admin user or not',
				active int(1) NOT NULL COMMENT 'Is active user or not',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE users ADD UNIQUE KEY email (email);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Commit all changes
		err = tx.Commit()
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
