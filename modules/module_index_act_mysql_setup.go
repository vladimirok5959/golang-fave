package modules

import (
	"os"
	"strconv"

	"golang-fave/engine/sqlw"
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
		db, err := sqlw.Open("mysql", pf_user+":"+pf_password+"@tcp("+pf_host+":"+pf_port+")/"+pf_name)
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

		// Table: blog_cat_post_rel
		if _, err = tx.Exec(
			`CREATE TABLE blog_cat_post_rel (
				post_id int(11) NOT NULL COMMENT 'Post id',
				category_id int(11) NOT NULL COMMENT 'Category id'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
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
				category int(11) NOT NULL,
				briefly text NOT NULL COMMENT 'Post brief content',
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

		// Table: notify_mail
		if _, err = tx.Exec(
			`CREATE TABLE notify_mail (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				email varchar(255) NOT NULL COMMENT 'Email address',
				subject varchar(800) NOT NULL COMMENT 'Email subject',
				message text NOT NULL COMMENT 'Email body',
				error text NOT NULL COMMENT 'Send error message',
				status int(1) NOT NULL COMMENT 'Sending status',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
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

		// Table: settings
		if _, err = tx.Exec(
			`CREATE TABLE settings (
				name varchar(255) NOT NULL COMMENT 'Setting name',
				value text NOT NULL COMMENT 'Setting value'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_cat_product_rel
		if _, err = tx.Exec(
			`CREATE TABLE shop_cat_product_rel (
				product_id int(11) NOT NULL COMMENT 'Product id',
				category_id int(11) NOT NULL COMMENT 'Category id'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_cats
		if _, err = tx.Exec(
			`CREATE TABLE shop_cats (
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

		// Table: shop_currencies
		if _, err = tx.Exec(
			`CREATE TABLE shop_currencies (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				name varchar(255) NOT NULL COMMENT 'Currency name',
				coefficient float(8,4) NOT NULL DEFAULT '1.0000' COMMENT 'Currency coefficient',
				code varchar(10) NOT NULL COMMENT 'Currency code',
				symbol varchar(5) NOT NULL COMMENT 'Currency symbol',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_filter_product_values
		if _, err = tx.Exec(
			`CREATE TABLE shop_filter_product_values (
				product_id int(11) NOT NULL COMMENT 'Product id',
				filter_value_id int(11) NOT NULL COMMENT 'Filter value id'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_filters
		if _, err = tx.Exec(
			`CREATE TABLE shop_filters (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				name varchar(255) NOT NULL COMMENT 'Filter name in CP',
				filter varchar(255) NOT NULL COMMENT 'Filter name in site',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_filters_values
		if _, err = tx.Exec(
			`CREATE TABLE shop_filters_values (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				filter_id int(11) NOT NULL COMMENT 'Filter id',
				name varchar(255) NOT NULL COMMENT 'Value name',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_product_images
		if _, err = tx.Exec(
			`CREATE TABLE shop_product_images (
				id int(11) NOT NULL AUTO_INCREMENT,
				product_id int(11) NOT NULL,
				filename varchar(255) NOT NULL,
				ord int(11) NOT NULL DEFAULT '0',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Table: shop_products
		if _, err = tx.Exec(
			`CREATE TABLE shop_products (
				id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
				parent_id int(11) DEFAULT NULL,
				user int(11) NOT NULL COMMENT 'User id',
				currency int(11) NOT NULL COMMENT 'Currency id',
				price float(8,2) NOT NULL COMMENT 'Product price',
				gname varchar(255) NOT NULL,
				name varchar(255) NOT NULL COMMENT 'Product name',
				alias varchar(255) NOT NULL COMMENT 'Product alias',
				vendor varchar(255) NOT NULL,
				quantity int(11) NOT NULL,
				category int(11) NOT NULL,
				briefly text NOT NULL COMMENT 'Product brief content',
				content text NOT NULL COMMENT 'Product content',
				datetime datetime NOT NULL COMMENT 'Creation date/time',
				active int(1) NOT NULL COMMENT 'Is active product or not',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		); err != nil {
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

		// Demo datas
		if _, err = tx.Exec(
			`INSERT INTO blog_cats (id, user, name, alias, lft, rgt)
				VALUES
			(1, 1, 'ROOT', 'ROOT', 1, 24),
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
			(12, 1, 'Film industry', 'film-industry', 19, 20);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO blog_cat_post_rel (post_id, category_id) VALUES (1, 9), (2, 12), (3, 8);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO blog_posts SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				category = ?,
				briefly = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			1,
			"Why should we eat wholesome food?",
			"why-should-we-eat-wholesome-food",
			9,
			"<p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO blog_posts SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				category = ?,
				briefly = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			2,
			1,
			"Latest top space movies",
			"latest-top-space-movies",
			12,
			"<p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO blog_posts SET
				id = ?,
				user = ?,
				name = ?,
				alias = ?,
				category = ?,
				briefly = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			3,
			1,
			"The best juices for a child",
			"the-best-juices-for-a-child",
			8,
			"<p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
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
		if _, err = tx.Exec(
			`INSERT INTO settings (name, value) VALUES ('database_version', '000000014');`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_cat_product_rel (product_id, category_id)
				VALUES
			(1, 3);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_cats (id, user, name, alias, lft, rgt)
				VALUES
			(1, 1, 'ROOT', 'ROOT', 1, 6),
			(2, 1, 'Electronics', 'electronics', 2, 5),
			(3, 1, 'Mobile phones', 'mobile-phones', 3, 4);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_currencies (id, name, coefficient, code, symbol)
				VALUES
			(1, 'US Dollar', 1.0000, 'USD', '$');`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_filter_product_values (product_id, filter_value_id)
				VALUES
			(1, 3),
			(1, 7),
			(1, 9),
			(1, 10),
			(1, 11);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_filters (id, name, filter)
				VALUES
			(1, 'Mobile phones manufacturer', 'Manufacturer'),
			(2, 'Mobile phones memory', 'Memory'),
			(3, 'Mobile phones communication standard', 'Communication standard');`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_filters_values (id, filter_id, name)
				VALUES
			(1, 1, 'Apple'),
			(2, 1, 'Asus'),
			(3, 1, 'Samsung'),
			(4, 2, '16 Gb'),
			(5, 2, '32 Gb'),
			(6, 2, '64 Gb'),
			(7, 2, '128 Gb'),
			(8, 2, '256 Gb'),
			(9, 3, '4G'),
			(10, 3, '2G'),
			(11, 3, '3G');`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO shop_products SET
				id = ?,
				user = ?,
				currency = ?,
				price = ?,
				gname = ?,
				name = ?,
				alias = ?,
				vendor = ?,
				quantity = ?,
				category = ?,
				briefly = ?,
				content = ?,
				datetime = ?,
				active = ?
			;`,
			1,
			1,
			1,
			1000.00,
			"",
			"Samsung Galaxy S10",
			"samsung-galaxy-s10",
			"Samsung",
			"1",
			"3",
			"<p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent.</p>",
			"<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat in ante metus dictum at tempor commodo ullamcorper a. Et malesuada fames ac turpis egestas sed tempus urna et. Euismod elementum nisi quis eleifend. Nisi porta lorem mollis aliquam ut porttitor. Ac turpis egestas maecenas pharetra convallis posuere. Nunc non blandit massa enim nec dui. Commodo elit at imperdiet dui accumsan sit amet nulla. Viverra accumsan in nisl nisi scelerisque. Dui nunc mattis enim ut tellus. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Faucibus ornare suspendisse sed nisi lacus. Nulla facilisi morbi tempus iaculis. Ut eu sem integer vitae justo eget magna fermentum iaculis. Ullamcorper sit amet risus nullam eget felis eget nunc. Volutpat sed cras ornare arcu dui vivamus. Eget magna fermentum iaculis eu non diam.</p><p>Arcu ac tortor dignissim convallis aenean et tortor. Vitae auctor eu augue ut lectus arcu. Ac turpis egestas integer eget aliquet nibh praesent. Interdum velit euismod in pellentesque massa placerat duis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt. Nisl rhoncus mattis rhoncus urna neque viverra justo. Odio ut enim blandit volutpat. Ac auctor augue mauris augue neque gravida. Ut lectus arcu bibendum at varius vel. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Dolor sit amet consectetur adipiscing elit duis tristique. Semper quis lectus nulla at volutpat diam ut. Sapien eget mi proin sed.</p>",
			utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
			1,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(
			`INSERT INTO users (id, first_name, last_name, email, password, admin, active) VALUES (1, 'First Name', 'Last Name', 'example@example.com', '23463b99b62a72f26ed677cc556c44e8', 1, 1);`,
		); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// Indexes
		if _, err = tx.Exec(`ALTER TABLE blog_cat_post_rel ADD UNIQUE KEY post_category (post_id,category_id) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cat_post_rel ADD KEY FK_blog_cat_post_rel_post_id (post_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_cat_post_rel ADD KEY FK_blog_cat_post_rel_category_id (category_id);`); err != nil {
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
		if _, err = tx.Exec(`ALTER TABLE blog_cats ADD KEY FK_blog_cats_user (user);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_posts ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_posts ADD KEY FK_blog_posts_user (user);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE blog_posts ADD KEY FK_blog_posts_category (category);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE notify_mail ADD KEY status (status);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE pages ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE pages ADD KEY alias_active (alias,active) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE pages ADD KEY FK_pages_user (user);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE settings ADD UNIQUE KEY name (name);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cat_product_rel ADD UNIQUE KEY product_category (product_id,category_id) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_product_id (product_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_category_id (category_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cats ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cats ADD KEY lft (lft), ADD KEY rgt (rgt);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_cats ADD KEY FK_shop_cats_user (user);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filter_product_values ADD UNIQUE KEY product_filter_value (product_id,filter_value_id) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_product_id (product_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_filter_value_id (filter_value_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filters ADD KEY name (name);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filters_values ADD KEY FK_shop_filters_values_filter_id (filter_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_filters_values ADD KEY name (name);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_product_images ADD UNIQUE KEY product_filename (product_id,filename) USING BTREE;`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_product_images ADD KEY FK_shop_product_images_product_id (product_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD UNIQUE KEY alias (alias);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD KEY FK_shop_products_user (user);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD KEY FK_shop_products_currency (currency);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD KEY FK_shop_products_category (category);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD KEY FK_shop_products_parent_id (parent_id);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE shop_products ADD KEY name (name);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`ALTER TABLE users ADD UNIQUE KEY email (email);`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}

		// References
		if _, err = tx.Exec(`
			ALTER TABLE blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_post_id
			FOREIGN KEY (post_id) REFERENCES blog_posts (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_category_id
			FOREIGN KEY (category_id) REFERENCES blog_cats (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE blog_cats ADD CONSTRAINT FK_blog_cats_user
			FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE blog_posts ADD CONSTRAINT FK_blog_posts_user
			FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE blog_posts ADD CONSTRAINT FK_blog_posts_category
			FOREIGN KEY (category) REFERENCES blog_cats (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE pages ADD CONSTRAINT FK_pages_user
			FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_product_id
			FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_category_id
			FOREIGN KEY (category_id) REFERENCES shop_cats (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_cats ADD CONSTRAINT FK_shop_cats_user
			FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_product_id
			FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_filter_value_id
			FOREIGN KEY (filter_value_id) REFERENCES shop_filters_values (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_filters_values ADD CONSTRAINT FK_shop_filters_values_filter_id
			FOREIGN KEY (filter_id) REFERENCES shop_filters (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_product_images ADD CONSTRAINT FK_shop_product_images_product_id
			FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_user
			FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_currency
			FOREIGN KEY (currency) REFERENCES shop_currencies (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_category
			FOREIGN KEY (category) REFERENCES shop_cats (id) ON DELETE RESTRICT;
		`); err != nil {
			tx.Rollback()
			wrap.MsgError(err.Error())
			return
		}
		if _, err = tx.Exec(`
			ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_parent_id
			FOREIGN KEY (parent_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
		`); err != nil {
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

		// Reset robots.txt file
		f, err := os.Create(wrap.DTemplate + string(os.PathSeparator) + "robots.txt")
		if err == nil {
			defer f.Close()
			if _, err = f.WriteString("User-agent: *\r\nDisallow: /\r\n"); err != nil {
				wrap.MsgError(err.Error())
				return
			}
		}

		// Create first config file
		if err := wrap.ConfigSave(); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.ResetCacheBlocks()

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}
