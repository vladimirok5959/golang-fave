package fetdata

import (
	"math"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type BlogPagination struct {
	Num     string
	Link    string
	Current bool
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

	bufferCats map[string][]*BlogCategory
}

func (this *Blog) init() {
	if this == nil {
		return
	}
	sql_nums := `
		SELECT
			COUNT(*)
		FROM
			blog_posts
		WHERE
			active = 1
		;
	`
	sql_rows := `
		SELECT
			id,
			user,
			name,
			alias,
			briefly,
			content,
			UNIX_TIMESTAMP(datetime) as datetime,
			active
		FROM
			blog_posts
		WHERE
			active = 1
		ORDER BY
			id DESC
		LIMIT ?, ?;
	`

	// Category selected
	if this.category != nil {
		var cat_ids []string
		if rows, err := this.wrap.DB.Query(
			`SELECT
				node.id
			FROM
				blog_cats AS node,
				blog_cats AS parent
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
				if err := rows.Scan(&cat_id); err == nil {
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
						blog_posts
						LEFT JOIN blog_cat_post_rel ON blog_cat_post_rel.post_id = blog_posts.id
					WHERE
						blog_posts.active = 1 AND
						blog_cat_post_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
					GROUP BY
						blog_posts.id
				) AS tbl
			;
		`
		sql_rows = `
			SELECT
				blog_posts.id,
				blog_posts.user,
				blog_posts.name,
				blog_posts.alias,
				blog_posts.briefly,
				blog_posts.content,
				UNIX_TIMESTAMP(blog_posts.datetime) AS datetime,
				blog_posts.active
			FROM
				blog_posts
				LEFT JOIN blog_cat_post_rel ON blog_cat_post_rel.post_id = blog_posts.id
			WHERE
				blog_posts.active = 1 AND
				blog_cat_post_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
			GROUP BY
				blog_posts.id
			ORDER BY
				blog_posts.id DESC
			LIMIT ?, ?;
		`
	}

	if err := this.wrap.DB.QueryRow(sql_nums).Scan(&this.postsCount); err == nil {
		// TODO: to control panel settings
		this.postsPerPage = 5
		this.postsMaxPage = int(math.Ceil(float64(this.postsCount) / float64(this.postsPerPage)))
		this.postsCurrPage = this.wrap.GetCurrentPage(this.postsMaxPage)
		offset := this.postsCurrPage*this.postsPerPage - this.postsPerPage
		if rows, err := this.wrap.DB.Query(sql_rows, offset, this.postsPerPage); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := utils.MySql_blog_post{}
				if err := rows.Scan(&row.A_id, &row.A_user, &row.A_name, &row.A_alias, &row.A_briefly, &row.A_content, &row.A_datetime, &row.A_active); err == nil {
					this.posts = append(this.posts, &BlogPost{object: &row})
				}
			}
		}
	}

	// Build pagination
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

func (this *Blog) Categories(mlvl int) []*BlogCategory {
	if this == nil {
		return []*BlogCategory{}
	}
	if this.bufferCats == nil {
		this.bufferCats = map[string][]*BlogCategory{}
	}
	key := ""
	where := ``
	if mlvl > 0 {
		where += `AND tbl.depth <= ` + utils.IntToStr(mlvl)
	}
	if _, ok := this.bufferCats[key]; !ok {
		var cats []*BlogCategory
		if rows, err := this.wrap.DB.Query(`
			SELECT
				tbl.*
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
						blog_cats AS node,
						blog_cats AS parent
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
					GROUP BY
						node.id
					ORDER BY
						node.lft ASC
				) AS tbl
			WHERE
				tbl.id > 1
				` + where + `
			;
		`); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := utils.MySql_blog_category{}
				var Depth int
				if err := rows.Scan(&row.A_id, &row.A_user, &row.A_name, &row.A_alias, &row.A_lft, &row.A_rgt, &Depth); err == nil {
					cats = append(cats, &BlogCategory{object: &row, depth: Depth})
				}
			}
		}
		this.bufferCats[key] = cats
	}
	return this.bufferCats[key]
}
