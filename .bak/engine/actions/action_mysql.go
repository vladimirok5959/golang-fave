package actions

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"strconv"

	utils "golang-fave/engine/wrapper/utils"
)

func (this *Action) Action_mysql() {
	pf_host := this.wrapper.R.FormValue("host")
	pf_port := this.wrapper.R.FormValue("port")
	pf_name := this.wrapper.R.FormValue("name")
	pf_user := this.wrapper.R.FormValue("user")
	pf_password := this.wrapper.R.FormValue("password")

	if pf_host == "" {
		this.msg_error(`Please specify host for MySQL connection`)
		return
	}

	if pf_port == "" {
		this.msg_error(`Please specify host port for MySQL connection`)
		return
	}

	if _, err := strconv.Atoi(pf_port); err != nil {
		this.msg_error(`MySQL host port must be integer number`)
		return
	}

	if pf_name == "" {
		this.msg_error(`Please specify MySQL database name`)
		return
	}

	if pf_user == "" {
		this.msg_error(`Please specify MySQL user`)
		return
	}

	// Try connect to mysql
	db, err := sql.Open("mysql", pf_user+":"+pf_password+"@tcp("+pf_host+":"+pf_port+")/"+pf_name)
	if err != nil {
		this.msg_error(err.Error())
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Try to install all tables
	_, err = db.Query(fmt.Sprintf(
		"CREATE TABLE `%s`.`users` (`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI', `first_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'User first name', `last_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'User last name', `email` VARCHAR(64) NOT NULL COMMENT 'User email', `password` VARCHAR(32) NOT NULL COMMENT 'User password (MD5)', PRIMARY KEY (`id`)) ENGINE = InnoDB;",
		pf_name))
	if err != nil {
		this.msg_error(err.Error())
		return
	}
	_, err = db.Query(fmt.Sprintf(
		"CREATE TABLE `%s`.`pages` (`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI', `user` int(11) NOT NULL COMMENT 'User id', `name` varchar(255) NOT NULL COMMENT 'Page name', `slug` varchar(255) NOT NULL COMMENT 'Page url part', `content` text NOT NULL COMMENT 'Page content', `meta_title` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta title', `meta_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta keywords', `meta_description` varchar(510) NOT NULL DEFAULT '' COMMENT 'Page meta description', `datetime` datetime NOT NULL COMMENT 'Creation date/time', `status` enum('draft','public','trash') NOT NULL COMMENT 'Page status', PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;",
		pf_name))
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Save mysql config file
	err = utils.MySqlConfigWrite(this.wrapper.DirVHostHome, pf_host, pf_port, pf_name, pf_user, pf_password)
	if err != nil {
		this.msg_error(err.Error())
		return
	}

	// Reload current page
	this.write(`window.location.reload(false);`)
}
