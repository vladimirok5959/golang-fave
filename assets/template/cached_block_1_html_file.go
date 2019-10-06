package template

var VarCachedBlock_1HtmlFile = []byte(`{{if not (eq $.Data.Module "404")}}
	{{if eq $.Data.Module "index"}}
		{{$.Data.Page.Name}}
	{{else if or (eq $.Data.Module "blog") (eq $.Data.Module "blog-post") (eq $.Data.Module "blog-category")}}
		{{if eq $.Data.Module "blog-category"}}
			Posts of category "{{$.Data.Blog.Category.Name}}" | Blog
		{{else if eq $.Data.Module "blog-post"}}
			{{$.Data.Blog.Post.Name}} | Blog
		{{else}}
			Latest posts | Blog
		{{end}}
	{{else if or (eq $.Data.Module "shop") (eq $.Data.Module "shop-product") (eq $.Data.Module "shop-category")}}
		{{if eq $.Data.Module "shop-category"}}
			Products of category "{{$.Data.Shop.Category.Name}}" | Shop
		{{else if eq $.Data.Module "shop-product"}}
			{{$.Data.Shop.Product.Name}} {{$.Data.Shop.Product.Id}} | Shop
		{{else}}
			Latest products | Shop
		{{end}}
	{{end}}
{{else}}
	Error 404
{{end}}`)
