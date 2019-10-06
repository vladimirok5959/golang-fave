package template

var VarCachedBlock_3HtmlFile = []byte(`<nav class="navbar navbar-expand-lg navbar-light bg-light navbar-cats">
	<div class="container">
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarSupportedContent">
			<ul class="navbar-nav mr-auto">
				{{range $.Data.Shop.Categories 0 1}}
					{{if .HaveChilds}}
						<li class="nav-item dropdown">
							<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">{{.Name}}</a>
							<div class="dropdown-menu" aria-labelledby="navbarDropdown">
								<a class="dropdown-item" href="{{.Permalink}}">All products</a>
								<div class="dropdown-divider"></div>
								{{range $index, $subcat := $.Data.Shop.Categories .Id 1}}
									<a class="dropdown-item" href="{{$subcat.Permalink}}">{{$subcat.Name}}</a>
								{{end}}
							</div>
						</li>
					{{else}}
						<li class="nav-item">
							<a class="nav-link" href="{{.Permalink}}">{{.Name}}</a>
						</li>
					{{end}}
				{{end}}
			</ul>
		</div>
	</div>
</nav>
{{if or (eq $.Data.Module "blog") (eq $.Data.Module "blog-category") (eq $.Data.Module "blog-post")}}
	<div class="container clear-top pt-4">
		<nav aria-label="breadcrumb">
			<ol class="breadcrumb mb-0">
				{{if eq $.Data.Module "blog"}}
					<li class="breadcrumb-item">Blog</li>
				{{else}}
					<li class="breadcrumb-item"><a href="/blog/">Blog</a></li>
				{{end}}
				{{if eq $.Data.Module "blog-category"}}
					{{if $.Data.Blog.Category.Parent.Parent.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Category.Parent.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Category.Parent.Parent.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Blog.Category.Parent.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Category.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Category.Parent.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Blog.Category.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Category.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Category.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Blog.Category.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Category.Parent.Parent.Permalink}}">{{$.Data.Blog.Category.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Blog.Category.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Category.Parent.Permalink}}">{{$.Data.Blog.Category.Parent.Name}}</a></li>
					{{end}}
					<li class="breadcrumb-item">{{$.Data.Blog.Category.Name}}</li>
				{{end}}
				{{if eq $.Data.Module "blog-post"}}
					{{if $.Data.Blog.Post.Category.Id}}
						{{if $.Data.Blog.Post.Category.Parent.Parent.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Blog.Post.Category.Parent.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Blog.Post.Category.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Permalink}}">{{$.Data.Blog.Post.Category.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Blog.Post.Category.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Parent.Parent.Permalink}}">{{$.Data.Blog.Post.Category.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Blog.Post.Category.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Parent.Permalink}}">{{$.Data.Blog.Post.Category.Parent.Name}}</a></li>
						{{end}}
						<li class="breadcrumb-item"><a href="{{$.Data.Blog.Post.Category.Permalink}}">{{$.Data.Blog.Post.Category.Name}}</a></li>
					{{end}}
					<li class="breadcrumb-item active">{{$.Data.Blog.Post.Name}}</li>
				{{end}}
			</ol>
		</nav>
	</div>
{{end}}
{{if or (eq $.Data.Module "shop") (eq $.Data.Module "shop-category") (eq $.Data.Module "shop-product")}}
	<div class="container clear-top pt-4">
		<nav aria-label="breadcrumb">
			<ol class="breadcrumb mb-0">
				{{if eq $.Data.Module "shop"}}
					<li class="breadcrumb-item">Shop</li>
				{{else}}
					<li class="breadcrumb-item"><a href="/shop/">Shop</a></li>
				{{end}}
				{{if eq $.Data.Module "shop-category"}}
					{{if $.Data.Shop.Category.Parent.Parent.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Category.Parent.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Category.Parent.Parent.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Shop.Category.Parent.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Category.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Category.Parent.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Shop.Category.Parent.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Category.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Category.Parent.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Shop.Category.Parent.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Category.Parent.Parent.Permalink}}">{{$.Data.Shop.Category.Parent.Parent.Name}}</a></li>
					{{end}}
					{{if $.Data.Shop.Category.Parent}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Category.Parent.Permalink}}">{{$.Data.Shop.Category.Parent.Name}}</a></li>
					{{end}}
					<li class="breadcrumb-item">{{$.Data.Shop.Category.Name}}</li>
				{{end}}
				{{if eq $.Data.Module "shop-product"}}
					{{if $.Data.Shop.Product.Category.Id}}
						{{if $.Data.Shop.Product.Category.Parent.Parent.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Shop.Product.Category.Parent.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Shop.Product.Category.Parent.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Permalink}}">{{$.Data.Shop.Product.Category.Parent.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Shop.Product.Category.Parent.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Parent.Parent.Permalink}}">{{$.Data.Shop.Product.Category.Parent.Parent.Name}}</a></li>
						{{end}}
						{{if $.Data.Shop.Product.Category.Parent}}
							<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Parent.Permalink}}">{{$.Data.Shop.Product.Category.Parent.Name}}</a></li>
						{{end}}
						<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Permalink}}">{{$.Data.Shop.Product.Category.Name}}</a></li>
					{{end}}
					<li class="breadcrumb-item active">{{$.Data.Shop.Product.Name}} ({{$.Data.Shop.Product.Id}})</li>
				{{end}}
			</ol>
		</nav>
	</div>
{{end}}`)
