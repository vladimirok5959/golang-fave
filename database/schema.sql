CREATE TABLE `blog_cats` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`user` int(11) NOT NULL COMMENT 'User id',
	`name` varchar(255) NOT NULL COMMENT 'Category name',
	`alias` varchar(255) NOT NULL COMMENT 'Category alias',
	`lft` int(11) NOT NULL COMMENT 'For nested set model',
	`rgt` int(11) NOT NULL COMMENT 'For nested set model',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `blog_cats` ADD UNIQUE KEY `alias` (`alias`);
ALTER TABLE `blog_cats` ADD KEY `lft` (`lft`), ADD KEY `rgt` (`rgt`);

CREATE TABLE `blog_cat_post_rel` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`post_id` int(11) NOT NULL COMMENT 'Post id',
	`category_id` int(11) NOT NULL COMMENT 'Category id',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `blog_cat_post_rel` ADD KEY `post_id` (`post_id`), ADD KEY `category_id` (`category_id`);

CREATE TABLE `blog_posts` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`user` int(11) NOT NULL COMMENT 'User id',
	`name` varchar(255) NOT NULL COMMENT 'Post name',
	`alias` varchar(255) NOT NULL COMMENT 'Post alias',
	`content` text NOT NULL COMMENT 'Post content',
	`datetime` datetime NOT NULL COMMENT 'Creation date/time',
	`active` int(1) NOT NULL COMMENT 'Is active post or not',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `blog_posts` ADD UNIQUE KEY `alias` (`alias`);

CREATE TABLE `pages` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`user` int(11) NOT NULL COMMENT 'User id',
	`name` varchar(255) NOT NULL COMMENT 'Page name',
	`alias` varchar(255) NOT NULL COMMENT 'Page url part',
	`content` text NOT NULL COMMENT 'Page content',
	`meta_title` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta title',
	`meta_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT 'Page meta keywords',
	`meta_description` varchar(510) NOT NULL DEFAULT '' COMMENT 'Page meta description',
	`datetime` datetime NOT NULL COMMENT 'Creation date/time',
	`active` int(1) NOT NULL COMMENT 'Is active page or not',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `pages` ADD UNIQUE KEY `alias` (`alias`);

CREATE TABLE `users` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`first_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User first name',
	`last_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'User last name',
	`email` varchar(64) NOT NULL COMMENT 'User email',
	`password` varchar(32) NOT NULL COMMENT 'User password (MD5)',
	`admin` int(1) NOT NULL COMMENT 'Is admin user or not',
	`active` int(1) NOT NULL COMMENT 'Is active user or not',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `users` ADD UNIQUE KEY `email` (`email`);