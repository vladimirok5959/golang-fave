package fetdata

import (
	"math"
	"sort"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type BlogPagination struct {
	Num     string
	Link    string
	Current bool
	Dots    bool
}

type Blog struct {
	wrap     *wrapper.Wrapper
	category *BlogCategory
	post     *BlogPost

	posts          []*BlogPost
	postsCount     int
	postsPerPage   int
	postsMaxPage   int
	postsCurrPage  int
	pagination     []*BlogPagination
	paginationPrev *BlogPagination
	paginationNext *BlogPagination

	bufferCats map[int]*utils.MySql_blog_category
}

func (this *Blog) load() *Blog {
	if this == nil {
		return this
	}
	sql_nums := `
		SELECT
			COUNT(*)
		FROM
			fave_blog_posts
		WHERE
			active = 1
		;
	`
	sql_rows := `
		SELECT
			fave_blog_posts.id,
			fave_blog_posts.user,
			fave_blog_posts.name,
			fave_blog_posts.alias,
			fave_blog_posts.category,
			fave_blog_posts.briefly,
			fave_blog_posts.content,
			UNIX_TIMESTAMP(fave_blog_posts.datetime) as datetime,
			fave_blog_posts.active,
			fave_users.id,
			fave_users.first_name,
			fave_users.last_name,
			fave_users.email,
			fave_users.admin,
			fave_users.active,
			cats.id,
			cats.user,
			cats.name,
			cats.alias,
			cats.lft,
			cats.rgt,
			cats.depth,
			cats.parent_id
		FROM
			fave_blog_posts
			LEFT JOIN fave_users ON fave_users.id = fave_blog_posts.user
			LEFT JOIN (
				SELECT
					main.id,
					main.user,
					main.name,
					main.alias,
					main.lft,
					main.rgt,
					main.depth,
					parent.id AS parent_id
				FROM
					(
						SELECT
							node.id,
							node.user,
							node.name,
							node.alias,
							node.lft,
							node.rgt,
							(COUNT(parent.id) - 1) AS depth
						FROM
							fave_blog_cats AS node,
							fave_blog_cats AS parent
						WHERE
							node.lft BETWEEN parent.lft AND parent.rgt
						GROUP BY
							node.id
						ORDER BY
							node.lft ASC
					) AS main
					LEFT JOIN (
						SELECT
							node.id,
							node.user,
							node.name,
							node.alias,
							node.lft,
							node.rgt,
							(COUNT(parent.id) - 0) AS depth
						FROM
							fave_blog_cats AS node,
							fave_blog_cats AS parent
						WHERE
							node.lft BETWEEN parent.lft AND parent.rgt
						GROUP BY
							node.id
						ORDER BY
							node.lft ASC
					) AS parent ON
					parent.depth = main.depth AND
					main.lft > parent.lft AND
					main.rgt < parent.rgt
				WHERE
					main.id > 1
				ORDER BY
					main.lft ASC
			) AS cats ON cats.id = fave_blog_posts.category
		WHERE
			fave_blog_posts.active = 1
		ORDER BY
			fave_blog_posts.id DESC
		LIMIT ?, ?;
	`

	// Category selected
	if this.category != nil {
		var cat_ids []string
		if rows, err := this.wrap.DB.Query(
			this.wrap.R.Context(),
			`SELECT
				node.id
			FROM
				fave_blog_cats AS node,
				fave_blog_cats AS parent
			WHERE
				node.lft BETWEEN parent.lft AND parent.rgt AND
				node.id > 1 AND
				parent.id = ?
			GROUP BY
				node.id
			ORDER BY
				node.lft ASC
			;`,
			this.category.Id(),
		); err == nil {
			defer rows.Close()
			for rows.Next() {
				var cat_id string
				if err := rows.Scan(&cat_id); *this.wrap.LogCpError(&err) == nil {
					cat_ids = append(cat_ids, cat_id)
				}
			}
		}
		sql_nums = `
			SELECT
				COUNT(*)
			FROM
				(
					SELECT
						COUNT(*)
					FROM
						fave_blog_posts
						LEFT JOIN fave_blog_cat_post_rel ON fave_blog_cat_post_rel.post_id = fave_blog_posts.id
					WHERE
						fave_blog_posts.active = 1 AND
						fave_blog_cat_post_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
					GROUP BY
						fave_blog_posts.id
				) AS tbl
			;
		`
		sql_rows = `
			SELECT
				fave_blog_posts.id,
				fave_blog_posts.user,
				fave_blog_posts.name,
				fave_blog_posts.alias,
				fave_blog_posts.category,
				fave_blog_posts.briefly,
				fave_blog_posts.content,
				UNIX_TIMESTAMP(fave_blog_posts.datetime) AS datetime,
				fave_blog_posts.active,
				fave_users.id,
				fave_users.first_name,
				fave_users.last_name,
				fave_users.email,
				fave_users.admin,
				fave_users.active,
				cats.id,
				cats.user,
				cats.name,
				cats.alias,
				cats.lft,
				cats.rgt,
				cats.depth,
				cats.parent_id
			FROM
				fave_blog_posts
				LEFT JOIN fave_blog_cat_post_rel ON fave_blog_cat_post_rel.post_id = fave_blog_posts.id
				LEFT JOIN fave_users ON fave_users.id = fave_blog_posts.user
				LEFT JOIN (
					SELECT
						main.id,
						main.user,
						main.name,
						main.alias,
						main.lft,
						main.rgt,
						main.depth,
						parent.id AS parent_id
					FROM
						(
							SELECT
								node.id,
								node.user,
								node.name,
								node.alias,
								node.lft,
								node.rgt,
								(COUNT(parent.id) - 1) AS depth
							FROM
								fave_blog_cats AS node,
								fave_blog_cats AS parent
							WHERE
								node.lft BETWEEN parent.lft AND parent.rgt
							GROUP BY
								node.id
							ORDER BY
								node.lft ASC
						) AS main
						LEFT JOIN (
							SELECT
								node.id,
								node.user,
								node.name,
								node.alias,
								node.lft,
								node.rgt,
								(COUNT(parent.id) - 0) AS depth
							FROM
								fave_blog_cats AS node,
								fave_blog_cats AS parent
							WHERE
								node.lft BETWEEN parent.lft AND parent.rgt
							GROUP BY
								node.id
							ORDER BY
								node.lft ASC
						) AS parent ON
						parent.depth = main.depth AND
						main.lft > parent.lft AND
						main.rgt < parent.rgt
					WHERE
						main.id > 1
					ORDER BY
						main.lft ASC
				) AS cats ON cats.id = fave_blog_posts.category
			WHERE
				fave_blog_posts.active = 1 AND
				fave_blog_cat_post_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
			GROUP BY
				fave_blog_posts.id,
				cats.parent_id
			ORDER BY
				fave_blog_posts.id DESC
			LIMIT ?, ?;
		`
	}

	if err := this.wrap.DB.QueryRow(this.wrap.R.Context(), sql_nums).Scan(&this.postsCount); *this.wrap.LogCpError(&err) == nil {
		if this.category == nil {
			this.postsPerPage = (*this.wrap.Config).Blog.Pagination.Index
		} else {
			this.postsPerPage = (*this.wrap.Config).Blog.Pagination.Category
		}
		this.postsMaxPage = int(math.Ceil(float64(this.postsCount) / float64(this.postsPerPage)))
		this.postsCurrPage = this.wrap.GetCurrentPage(this.postsMaxPage)
		offset := this.postsCurrPage*this.postsPerPage - this.postsPerPage
		if rows, err := this.wrap.DB.Query(this.wrap.R.Context(), sql_rows, offset, this.postsPerPage); err == nil {
			defer rows.Close()
			for rows.Next() {
				rp := utils.MySql_blog_post{}
				ru := utils.MySql_user{}
				ro := utils.MySql_blog_category{}
				if err := rows.Scan(
					&rp.A_id,
					&rp.A_user,
					&rp.A_name,
					&rp.A_alias,
					&rp.A_category,
					&rp.A_briefly,
					&rp.A_content,
					&rp.A_datetime,
					&rp.A_active,
					&ru.A_id,
					&ru.A_first_name,
					&ru.A_last_name,
					&ru.A_email,
					&ru.A_admin,
					&ru.A_active,
					&ro.A_id,
					&ro.A_user,
					&ro.A_name,
					&ro.A_alias,
					&ro.A_lft,
					&ro.A_rgt,
					&ro.A_depth,
					&ro.A_parent,
				); *this.wrap.LogCpError(&err) == nil {
					this.posts = append(this.posts, &BlogPost{
						wrap:     this.wrap,
						object:   &rp,
						user:     &User{wrap: this.wrap, object: &ru},
						category: &BlogCategory{wrap: this.wrap, object: &ro},
					})
				}
			}
		}
	}

	// Build pagination
	if true {
		for i := 1; i < this.postsCurrPage; i++ {
			if this.postsCurrPage >= 5 && i > 1 && i < this.postsCurrPage-1 {
				continue
			}
			if this.postsCurrPage >= 5 && i > 1 && i < this.postsCurrPage {
				this.pagination = append(this.pagination, &BlogPagination{
					Dots: true,
				})
			}
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &BlogPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: false,
			})
		}

		// Current page
		link := this.wrap.R.URL.Path
		if this.postsCurrPage > 1 {
			link = link + "?p=" + utils.IntToStr(this.postsCurrPage)
		}
		this.pagination = append(this.pagination, &BlogPagination{
			Num:     utils.IntToStr(this.postsCurrPage),
			Link:    link,
			Current: true,
		})

		for i := this.postsCurrPage + 1; i <= this.postsMaxPage; i++ {
			if this.postsCurrPage < this.postsMaxPage-3 && i == this.postsCurrPage+3 {
				this.pagination = append(this.pagination, &BlogPagination{
					Dots: true,
				})
			}
			if this.postsCurrPage < this.postsMaxPage-3 && i > this.postsCurrPage+1 && i <= this.postsMaxPage-1 {
				continue
			}
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &BlogPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: false,
			})
		}
	} else {
		for i := 1; i <= this.postsMaxPage; i++ {
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &BlogPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: i == this.postsCurrPage,
			})
		}
	}

	// Pagination prev/next
	if this.postsMaxPage > 1 {
		link := this.wrap.R.URL.Path
		if this.postsCurrPage-1 > 1 {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.postsCurrPage-1)
		}
		this.paginationPrev = &BlogPagination{
			Num:     utils.IntToStr(this.postsCurrPage - 1),
			Link:    link,
			Current: this.postsCurrPage <= 1,
		}
		if this.postsCurrPage >= 1 && this.postsCurrPage < this.postsMaxPage {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.postsCurrPage+1)
		} else {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.postsMaxPage)
		}
		this.paginationNext = &BlogPagination{
			Num:     utils.IntToStr(this.postsCurrPage + 1),
			Link:    link,
			Current: this.postsCurrPage >= this.postsMaxPage,
		}
	}

	return this
}

func (this *Blog) preload_cats() {
	if this.bufferCats == nil {
		this.bufferCats = map[int]*utils.MySql_blog_category{}
		if rows, err := this.wrap.DB.Query(
			this.wrap.R.Context(),
			`SELECT
				main.id,
				main.user,
				main.name,
				main.alias,
				main.lft,
				main.rgt,
				main.depth,
				parent.id AS parent_id
			FROM
				(
					SELECT
						node.id,
						node.user,
						node.name,
						node.alias,
						node.lft,
						node.rgt,
						(COUNT(parent.id) - 1) AS depth
					FROM
						fave_blog_cats AS node,
						fave_blog_cats AS parent
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
					GROUP BY
						node.id
					ORDER BY
						node.lft ASC
				) AS main
				LEFT JOIN (
					SELECT
						node.id,
						node.user,
						node.name,
						node.alias,
						node.lft,
						node.rgt,
						(COUNT(parent.id) - 0) AS depth
					FROM
						fave_blog_cats AS node,
						fave_blog_cats AS parent
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
					GROUP BY
						node.id
					ORDER BY
						node.lft ASC
				) AS parent ON
				parent.depth = main.depth AND
				main.lft > parent.lft AND
				main.rgt < parent.rgt
			WHERE
				main.id > 1
			ORDER BY
				main.lft ASC
			;
		`); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := utils.MySql_blog_category{}
				if err := rows.Scan(
					&row.A_id,
					&row.A_user,
					&row.A_name,
					&row.A_alias,
					&row.A_lft,
					&row.A_rgt,
					&row.A_depth,
					&row.A_parent,
				); *this.wrap.LogCpError(&err) == nil {
					this.bufferCats[row.A_id] = &row
					if _, ok := this.bufferCats[row.A_parent]; ok {
						this.bufferCats[row.A_parent].A_childs = true
					}
				}
			}
		}
	}
}

func (this *Blog) Category() *BlogCategory {
	if this == nil {
		return nil
	}
	return this.category
}

func (this *Blog) Post() *BlogPost {
	if this == nil {
		return nil
	}
	return this.post
}

func (this *Blog) HavePosts() bool {
	if this == nil {
		return false
	}
	if len(this.posts) <= 0 {
		return false
	}
	return true
}

func (this *Blog) Posts() []*BlogPost {
	if this == nil {
		return []*BlogPost{}
	}
	return this.posts
}

func (this *Blog) PostsCount() int {
	if this == nil {
		return 0
	}
	return this.postsCount
}

func (this *Blog) PostsPerPage() int {
	if this == nil {
		return 0
	}
	return this.postsPerPage
}

func (this *Blog) PostsMaxPage() int {
	if this == nil {
		return 0
	}
	return this.postsMaxPage
}

func (this *Blog) PostsCurrPage() int {
	if this == nil {
		return 0
	}
	return this.postsCurrPage
}

func (this *Blog) Pagination() []*BlogPagination {
	if this == nil {
		return []*BlogPagination{}
	}
	return this.pagination
}

func (this *Blog) PaginationPrev() *BlogPagination {
	if this == nil {
		return nil
	}
	return this.paginationPrev
}

func (this *Blog) PaginationNext() *BlogPagination {
	if this == nil {
		return nil
	}
	return this.paginationNext
}

func (this *Blog) Categories(parent, depth int) []*BlogCategory {
	this.preload_cats()

	depth_tmp := 0
	result := []*BlogCategory{}

	for _, cat := range this.bufferCats {
		if parent <= 1 {
			if depth <= 0 {
				result = append(result, (&BlogCategory{wrap: this.wrap, object: cat}).load(&this.bufferCats))
			} else {
				if cat.A_depth <= depth {
					result = append(result, (&BlogCategory{wrap: this.wrap, object: cat}).load(&this.bufferCats))
				}
			}
		} else {
			if cat.A_parent == parent {
				if depth_tmp == 0 {
					depth_tmp = cat.A_depth
				}
				if depth <= 0 {
					result = append(result, (&BlogCategory{wrap: this.wrap, object: cat}).load(&this.bufferCats))
				} else {
					if (cat.A_depth - depth_tmp + 1) <= depth {
						result = append(result, (&BlogCategory{wrap: this.wrap, object: cat}).load(&this.bufferCats))
					}
				}
			}
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Left() < result[j].Left() })

	return result
}
