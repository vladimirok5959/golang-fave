package template

var VarCachedBlock_2HtmlFile = []byte(`{{if not (eq $.Data.Module "404")}}
	{{if eq $.Data.Module "index"}}
		{{if eq $.Data.Page.Alias "/"}}
			Welcome to home page
		{{else}}
			Welcome to some another page
		{{end}}
	{{else if or (eq $.Data.Module "blog") (eq $.Data.Module "blog-post") (eq $.Data.Module "blog-category")}}
		{{if eq $.Data.Module "blog-category"}}
			Blog category
		{{else if eq $.Data.Module "blog-post"}}
			Blog post
		{{else}}
			Blog
		{{end}}
	{{else if or (eq $.Data.Module "shop") (eq $.Data.Module "shop-product") (eq $.Data.Module "shop-category")}}
		{{if eq $.Data.Module "shop-category"}}
			Shop category
		{{else if eq $.Data.Module "shop-product"}}
			Shop product
		{{else}}
			Shop
		{{end}}
	{{end}}
{{else}}
	Oops, page is not found...
{{end}}`)
