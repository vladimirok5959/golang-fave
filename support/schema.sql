# Tables with keys
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
ALTER TABLE `blog_cats` ADD KEY `FK_blog_cats_user` (`user`);

CREATE TABLE `blog_cat_post_rel` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`post_id` int(11) NOT NULL COMMENT 'Post id',
	`category_id` int(11) NOT NULL COMMENT 'Category id',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `blog_cat_post_rel` ADD KEY `post_id` (`post_id`), ADD KEY `category_id` (`category_id`);
ALTER TABLE `blog_cat_post_rel` ADD UNIQUE KEY `post_category` (`post_id`,`category_id`) USING BTREE;
ALTER TABLE `blog_cat_post_rel` ADD KEY `FK_blog_cat_post_rel_post_id` (`post_id`);
ALTER TABLE `blog_cat_post_rel` ADD KEY `FK_blog_cat_post_rel_category_id` (`category_id`);

CREATE TABLE `blog_posts` (
	`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	`user` int(11) NOT NULL COMMENT 'User id',
	`name` varchar(255) NOT NULL COMMENT 'Post name',
	`alias` varchar(255) NOT NULL COMMENT 'Post alias',
	`briefly` text NOT NULL COMMENT 'Post brief content',
	`content` text NOT NULL COMMENT 'Post content',
	`datetime` datetime NOT NULL COMMENT 'Creation date/time',
	`active` int(1) NOT NULL COMMENT 'Is active post or not',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `blog_posts` ADD UNIQUE KEY `alias` (`alias`);
ALTER TABLE `blog_posts` ADD KEY `FK_blog_posts_user` (`user`);

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
ALTER TABLE `pages` ADD KEY `alias_active` (`alias`,`active`) USING BTREE;
ALTER TABLE `pages` ADD KEY `FK_pages_user` (`user`);

CREATE TABLE `settings` (
	`name` varchar(255) NOT NULL COMMENT 'Setting name',
	`value` text NOT NULL COMMENT 'Setting value',
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `settings` ADD UNIQUE KEY `name` (`name`);

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

# References
ALTER TABLE `blog_cats` ADD CONSTRAINT `FK_blog_cats_user` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `blog_cat_post_rel` ADD CONSTRAINT `FK_blog_cat_post_rel_category_id` FOREIGN KEY (`category_id`) REFERENCES `blog_cats` (`id`) ON DELETE RESTRICT;
ALTER TABLE `blog_cat_post_rel` ADD CONSTRAINT `FK_blog_cat_post_rel_post_id` FOREIGN KEY (`post_id`) REFERENCES `blog_posts` (`id`) ON DELETE RESTRICT;
ALTER TABLE `blog_posts` ADD CONSTRAINT `FK_blog_posts_user` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `pages` ADD CONSTRAINT `FK_pages_user` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
