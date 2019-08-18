package template

var VarHeaderHtmlFile = []byte(`<!doctype html>
<html lang="en">
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="theme-color" content="#205081" />
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="{{$.System.PathCssBootstrap}}">

		<title>
			{{if not (eq $.Data.Module "404")}}
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
						{{$.Data.Shop.Product.Name}} | Shop
					{{else}}
						Latest products | Shop
					{{end}}
				{{end}}
			{{else}}
				Error 404
			{{end}}
		</title>
		<meta name="keywords" content="{{$.Data.Page.MetaKeywords}}" />
		<meta name="description" content="{{$.Data.Page.MetaDescription}}" />
		<link rel="shortcut icon" href="{{$.System.PathIcoFav}}" type="image/x-icon" />

		<!-- Template CSS file from template folder -->
		<link rel="stylesheet" href="{{$.System.PathThemeStyles}}?v=2">
	</head>
	<body class="fixed-top-bar">
		<div id="wrap">
			<nav class="navbar navbar-expand-lg navbar-light bg-light">
				<div class="container">
					<a class="navbar-brand" href="/">Fave {{$.System.InfoVersion}}</a>
					<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive" aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
						<span class="navbar-toggler-icon"></span>
					</button>
					<div class="collapse navbar-collapse" id="navbarResponsive">
						<ul class="navbar-nav ml-auto">
							<li class="nav-item{{if eq $.Data.Page.Alias "/"}} active{{end}}">
								<a class="nav-link" href="/">Home</a>
							</li>
							<li class="nav-item">
								<a class="nav-link{{if eq $.Data.Page.Alias "/another/"}} active{{end}}" href="/another/">Another</a>
							</li>
							<li class="nav-item">
								<a class="nav-link{{if eq $.Data.Page.Alias "/about/"}} active{{end}}" href="/about/">About</a>
							</li>
							<li class="nav-item">
								<a class="nav-link{{if or (eq $.Data.Module "blog") (eq $.Data.Module "blog-post") (eq $.Data.Module "blog-category")}} active{{end}}" href="/blog/">Blog</a>
							</li>
							<li class="nav-item">
								<a class="nav-link{{if or (eq $.Data.Module "shop") (eq $.Data.Module "shop-product") (eq $.Data.Module "shop-category")}} active{{end}}" href="/shop/">Shop</a>
							</li>
							<li class="nav-item">
								<a class="nav-link{{if eq $.Data.Module "404"}} active{{end}}" href="/not-existent-page/">404</a>
							</li>
						</ul>
					</div>
				</div>
			</nav>
			<div id="main">
				<div class="bg-fave">
					<div class="container">
						<h1 class="text-left text-white m-0 p-0 py-5">
							{{if not (eq $.Data.Module "404")}}
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
							{{end}}
						</h1>
					</div>
				</div>
				<div class="container clear-top">
					<div class="row pt-4">
						{{if or (eq $.Data.Module "shop") (eq $.Data.Module "shop-category")}}
							<div class="col-sm-5 col-md-4 col-lg-3">
								{{template "sidebar-left.html" .}}
							</div>
						{{end}}
						{{if or (eq $.Data.Module "shop-product")}}
							<div class="col-md-12">
						{{else}}
							<div class="col-sm-7 col-md-8 col-lg-9">
						{{end}}`)
