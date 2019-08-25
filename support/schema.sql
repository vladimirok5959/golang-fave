# Tables
CREATE TABLE blog_cats (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	user int(11) NOT NULL COMMENT 'User id',
	name varchar(255) NOT NULL COMMENT 'Category name',
	alias varchar(255) NOT NULL COMMENT 'Category alias',
	lft int(11) NOT NULL COMMENT 'For nested set model',
	rgt int(11) NOT NULL COMMENT 'For nested set model',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE blog_cat_post_rel (
	post_id int(11) NOT NULL COMMENT 'Post id',
	category_id int(11) NOT NULL COMMENT 'Category id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE blog_posts (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE pages (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE settings (
	name varchar(255) NOT NULL COMMENT 'Setting name',
	value text NOT NULL COMMENT 'Setting value'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_cat_product_rel (
	product_id int(11) NOT NULL COMMENT 'Product id',
	category_id int(11) NOT NULL COMMENT 'Category id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_cats (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	user int(11) NOT NULL COMMENT 'User id',
	name varchar(255) NOT NULL COMMENT 'Category name',
	alias varchar(255) NOT NULL COMMENT 'Category alias',
	lft int(11) NOT NULL COMMENT 'For nested set model',
	rgt int(11) NOT NULL COMMENT 'For nested set model',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_currencies (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	name varchar(255) NOT NULL COMMENT 'Currency name',
	coefficient float(8,4) NOT NULL DEFAULT '1.0000' COMMENT 'Currency coefficient',
	code varchar(10) NOT NULL COMMENT 'Currency code',
	symbol varchar(5) NOT NULL COMMENT 'Currency symbol',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_filter_product_values (
	product_id int(11) NOT NULL COMMENT 'Product id',
	filter_value_id int(11) NOT NULL COMMENT 'Filter value id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_filters (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	name varchar(255) NOT NULL COMMENT 'Filter name in CP',
	filter varchar(255) NOT NULL COMMENT 'Filter name in site',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_filters_values (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	filter_id int(11) NOT NULL COMMENT 'Filter id',
	name varchar(255) NOT NULL COMMENT 'Value name',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_product_images (
	product_id int(11) NOT NULL,
	filename varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE shop_products (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	user int(11) NOT NULL COMMENT 'User id',
	currency int(11) NOT NULL COMMENT 'Currency id',
	price float(8,2) NOT NULL COMMENT 'Product price',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE users (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
	first_name varchar(64) NOT NULL DEFAULT '' COMMENT 'User first name',
	last_name varchar(64) NOT NULL DEFAULT '' COMMENT 'User last name',
	email varchar(64) NOT NULL COMMENT 'User email',
	password varchar(32) NOT NULL COMMENT 'User password (MD5)',
	admin int(1) NOT NULL COMMENT 'Is admin user or not',
	active int(1) NOT NULL COMMENT 'Is active user or not',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Indexes
ALTER TABLE blog_cat_post_rel ADD UNIQUE KEY post_category (post_id,category_id) USING BTREE;
ALTER TABLE blog_cat_post_rel ADD KEY FK_blog_cat_post_rel_post_id (post_id);
ALTER TABLE blog_cat_post_rel ADD KEY FK_blog_cat_post_rel_category_id (category_id);
ALTER TABLE blog_cats ADD UNIQUE KEY alias (alias);
ALTER TABLE blog_cats ADD KEY lft (lft), ADD KEY rgt (rgt);
ALTER TABLE blog_cats ADD KEY FK_blog_cats_user (user);
ALTER TABLE blog_posts ADD UNIQUE KEY alias (alias);
ALTER TABLE blog_posts ADD KEY FK_blog_posts_user (user);
ALTER TABLE blog_posts ADD KEY FK_blog_posts_category (category);
ALTER TABLE pages ADD UNIQUE KEY alias (alias);
ALTER TABLE pages ADD KEY alias_active (alias,active) USING BTREE;
ALTER TABLE pages ADD KEY FK_pages_user (user);
ALTER TABLE settings ADD UNIQUE KEY name (name);
ALTER TABLE shop_cat_product_rel ADD UNIQUE KEY product_category (product_id,category_id) USING BTREE;
ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_product_id (product_id);
ALTER TABLE shop_cat_product_rel ADD KEY FK_shop_cat_product_rel_category_id (category_id);
ALTER TABLE shop_cats ADD UNIQUE KEY alias (alias);
ALTER TABLE shop_cats ADD KEY lft (lft), ADD KEY rgt (rgt);
ALTER TABLE shop_cats ADD KEY FK_shop_cats_user (user);
ALTER TABLE shop_filter_product_values ADD UNIQUE KEY product_filter_value (product_id,filter_value_id) USING BTREE;
ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_product_id (product_id);
ALTER TABLE shop_filter_product_values ADD KEY FK_shop_filter_product_values_filter_value_id (filter_value_id);
ALTER TABLE shop_filters ADD KEY name (name);
ALTER TABLE shop_filters_values ADD KEY FK_shop_filters_values_filter_id (filter_id);
ALTER TABLE shop_filters_values ADD KEY name (name);
ALTER TABLE shop_product_images ADD UNIQUE KEY product_filename (product_id,filename) USING BTREE;
ALTER TABLE shop_product_images ADD KEY FK_shop_product_images_product_id (product_id);
ALTER TABLE shop_products ADD UNIQUE KEY alias (alias);
ALTER TABLE shop_products ADD KEY FK_shop_products_user (user);
ALTER TABLE shop_products ADD KEY FK_shop_products_currency (currency);
ALTER TABLE shop_products ADD KEY FK_shop_products_category (category);
ALTER TABLE users ADD UNIQUE KEY email (email);

# References
ALTER TABLE blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_post_id FOREIGN KEY (post_id) REFERENCES blog_posts (id) ON DELETE RESTRICT;
ALTER TABLE blog_cat_post_rel ADD CONSTRAINT FK_blog_cat_post_rel_category_id FOREIGN KEY (category_id) REFERENCES blog_cats (id) ON DELETE RESTRICT;
ALTER TABLE blog_cats ADD CONSTRAINT FK_blog_cats_user FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
ALTER TABLE blog_posts ADD CONSTRAINT FK_blog_posts_user FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
ALTER TABLE blog_posts ADD CONSTRAINT FK_blog_posts_category FOREIGN KEY (category) REFERENCES blog_cats (id) ON DELETE RESTRICT;
ALTER TABLE pages ADD CONSTRAINT FK_pages_user FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_product_id FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
ALTER TABLE shop_cat_product_rel ADD CONSTRAINT FK_shop_cat_product_rel_category_id FOREIGN KEY (category_id) REFERENCES shop_cats (id) ON DELETE RESTRICT;
ALTER TABLE shop_cats ADD CONSTRAINT FK_shop_cats_user FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_product_id FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
ALTER TABLE shop_filter_product_values ADD CONSTRAINT FK_shop_filter_product_values_filter_value_id FOREIGN KEY (filter_value_id) REFERENCES shop_filters_values (id) ON DELETE RESTRICT;
ALTER TABLE shop_filters_values ADD CONSTRAINT FK_shop_filters_values_filter_id FOREIGN KEY (filter_id) REFERENCES shop_filters (id) ON DELETE RESTRICT;
ALTER TABLE shop_product_images ADD CONSTRAINT FK_shop_product_images_product_id FOREIGN KEY (product_id) REFERENCES shop_products (id) ON DELETE RESTRICT;
ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_user FOREIGN KEY (user) REFERENCES users (id) ON DELETE RESTRICT;
ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_currency FOREIGN KEY (currency) REFERENCES shop_currencies (id) ON DELETE RESTRICT;
ALTER TABLE shop_products ADD CONSTRAINT FK_shop_products_category FOREIGN KEY (category) REFERENCES shop_cats (id) ON DELETE RESTRICT;
