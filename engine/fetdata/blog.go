package fetdata

// import (
// 	"math"
// 	"strconv"
// 	"strings"

// 	"golang-fave/engine/sqlw"
// 	"golang-fave/utils"
// )

// func (this *FERData) postsGetCount(buf string, cat int) (int, int) {
// 	if cat == 0 {
// 		var num int
// 		if err := this.wrap.DB.QueryRow(`
// 			SELECT
// 				COUNT(*)
// 			FROM
// 				blog_posts
// 			WHERE
// 				active = 1
// 			;
// 		`).Scan(&num); err == nil {
// 			pear_page := 2
// 			max_pages := int(math.Ceil(float64(num) / float64(pear_page)))
// 			curr_page := 1
// 			p := this.wrap.R.URL.Query().Get("p")
// 			if p != "" {
// 				pi, err := strconv.Atoi(p)
// 				if err != nil {
// 					curr_page = 1
// 				} else {
// 					if pi < 1 {
// 						curr_page = 1
// 					} else if pi > max_pages {
// 						curr_page = max_pages
// 					} else {
// 						curr_page = pi
// 					}
// 				}
// 			}
// 			limit_offset := curr_page*pear_page - pear_page
// 			return limit_offset, pear_page
// 		}
// 	} else {
// 		var num int
// 		if err := this.wrap.DB.QueryRow(`
// 			SELECT
// 				COUNT(blog_posts.id)
// 			FROM
// 				blog_posts
// 				LEFT JOIN blog_cat_post_rel ON blog_cat_post_rel.post_id = blog_posts.id
// 			WHERE
// 				blog_posts.active = 1 AND
// 				blog_cat_post_rel.category_id = ?
// 			;
// 		`, cat).Scan(&num); err == nil {
// 			pear_page := 2
// 			max_pages := int(math.Ceil(float64(num) / float64(pear_page)))
// 			curr_page := 1
// 			p := this.wrap.R.URL.Query().Get("p")
// 			if p != "" {
// 				pi, err := strconv.Atoi(p)
// 				if err != nil {
// 					curr_page = 1
// 				} else {
// 					if pi < 1 {
// 						curr_page = 1
// 					} else if pi > max_pages {
// 						curr_page = max_pages
// 					} else {
// 						curr_page = pi
// 					}
// 				}
// 			}
// 			limit_offset := curr_page*pear_page - pear_page
// 			return limit_offset, pear_page
// 		}
// 	}
// 	return 0, 0
// }

// func (this *FERData) postsToBuffer(buf string, cat int, order string) {
// 	if this.bufferPosts == nil {
// 		this.bufferPosts = map[string][]*BlogPost{}
// 	}
// 	if _, ok := this.bufferPosts[buf]; !ok {
// 		var posts []*BlogPost

// 		limit_offset, pear_page := this.postsGetCount(buf, cat)

// 		var rows *sqlw.Rows
// 		var err error

// 		if cat == 0 {
// 			rows, err = this.wrap.DB.Query(`
// 				SELECT
// 					blog_posts.id,
// 					blog_posts.user,
// 					blog_posts.name,
// 					blog_posts.alias,
// 					blog_posts.content,
// 					UNIX_TIMESTAMP(blog_posts.datetime) AS datetime,
// 					blog_posts.active
// 				FROM
// 					blog_posts
// 				WHERE
// 					blog_posts.active = 1
// 				ORDER BY
// 					blog_posts.id `+order+`
// 				LIMIT ?, ?;
// 			`, limit_offset, pear_page)
// 		} else {
// 			rows, err = this.wrap.DB.Query(`
// 				SELECT
// 					blog_posts.id,
// 					blog_posts.user,
// 					blog_posts.name,
// 					blog_posts.alias,
// 					blog_posts.content,
// 					UNIX_TIMESTAMP(blog_posts.datetime) AS datetime,
// 					blog_posts.active
// 				FROM
// 					blog_posts
// 					LEFT JOIN blog_cat_post_rel ON blog_cat_post_rel.post_id = blog_posts.id
// 				WHERE
// 					blog_posts.active = 1 AND
// 					blog_cat_post_rel.category_id = ?
// 				ORDER BY
// 					blog_posts.id `+order+`
// 				LIMIT ?, ?;
// 			`, cat, limit_offset, pear_page)
// 		}

// 		if err == nil {
// 			var f_id int
// 			var f_user int
// 			var f_name string
// 			var f_alias string
// 			var f_content string
// 			var f_datetime int
// 			var f_active int
// 			for rows.Next() {
// 				err = rows.Scan(&f_id, &f_user, &f_name, &f_alias, &f_content, &f_datetime, &f_active)
// 				if err == nil {
// 					posts = append(posts, &BlogPost{
// 						id:       f_id,
// 						user:     f_user,
// 						name:     f_name,
// 						alias:    f_alias,
// 						content:  f_content,
// 						datetime: f_datetime,
// 						active:   f_active,
// 					})
// 				}
// 			}
// 			rows.Close()
// 		}
// 		this.bufferPosts[buf] = posts
// 	}
// }

// func (this *FERData) BlogPosts() []*BlogPost {
// 	return this.BlogPostsOrder("DESC")
// }

// func (this *FERData) BlogPostsOrder(order string) []*BlogPost {
// 	posts_order := "DESC"

// 	if strings.ToLower(order) == "asc" {
// 		posts_order = "ASC"
// 	}

// 	buf := "posts_" + posts_order
// 	this.postsToBuffer(buf, 0, posts_order)
// 	return this.bufferPosts[buf]
// }

// func (this *FERData) BlogPostsOfCat(cat int) []*BlogPost {
// 	return this.BlogPostsOfCatOrder(cat, "DESC")
// }

// func (this *FERData) BlogPostsOfCatOrder(cat int, order string) []*BlogPost {
// 	posts_order := "DESC"

// 	if strings.ToLower(order) == "asc" {
// 		posts_order = "ASC"
// 	}

// 	buf := "posts_" + posts_order + "_" + utils.IntToStr(cat)
// 	this.postsToBuffer(buf, cat, posts_order)
// 	return this.bufferPosts[buf]
// }
