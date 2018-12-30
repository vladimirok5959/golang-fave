-- version 0.0.0

CREATE TABLE `pages` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`parent` int(11) NOT NULL DEFAULT '0' COMMENT 'Parent page id',
	`user` int(11) NOT NULL COMMENT 'User id',
	`name` varchar(255) NOT NULL COMMENT 'Page name',
	`slug` varchar(255) NOT NULL COMMENT 'Page url part',
	`content` text NOT NULL COMMENT 'Page content',
	`meta_title` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta title',
	`meta_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta keywords',
	`meta_description` varchar(510) NOT NULL DEFAULT '' COMMENT 'Page meta description',
	`datetime` datetime NOT NULL COMMENT 'Creation date/time',
	`status` enum('draft','public','trash') NOT NULL COMMENT 'Page status',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`first_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User first name',
	`last_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User last name',
	`email` varchar(64) NOT NULL COMMENT 'User email',
	`password` varchar(32) NOT NULL COMMENT 'User password (MD5)',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
