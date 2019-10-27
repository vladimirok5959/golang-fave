package template

var VarShopHtmlFile = []byte(`{{template "header.html" .}}
<div class="mb-4">
	{{if $.Data.Shop.HaveProducts}}
		<div class="grid-products">
			{{range $.Data.Shop.Products}}
				<div class="card card-product">
					<div class="card-img-link">
						<a href="{{.Permalink}}">
							{{if .HaveImages }}
								<img class="card-img-top" src="{{.Image.Thumbnail1}}" alt="{{$.Data.EscapeString .Name}}">
							{{else}}
								<img class="card-img-top" src="{{$.Data.ImagePlaceholderHref}}" alt="{{$.Data.EscapeString .Name}}">
							{{end}}
						</a>
					</div>
					<div class="card-body">
						<h5 class="card-title">
							<a href="{{.Permalink}}">
								{{if ne .Group ""}}
									{{.Group}}
								{{else}}
									{{.Name}}
								{{end}}
							</a>
						</h5>
						<div class="card-text">{{.Briefly}}</div>
					</div>
					<div class="card-footer">
						<span class="price">{{.PriceNice}} {{$.Data.Shop.CurrentCurrency.Code}}</span><a href="{{.Permalink}}" class="btn btn-success">View</a>
					</div>
				</div>
			{{end}}
		</div>
	{{else}}
		<div class="card">
			<div class="card-body">
				Sorry, no products matched your criteria
			</div>
		</div>
	{{end}}
</div>
{{if $.Data.Shop.HaveProducts}}
	{{if gt $.Data.Shop.ProductsMaxPage 1 }}
		<nav>
			<ul class="pagination mb-4">
				{{if $.Data.Shop.PaginationPrev}}
					<li class="page-item{{if $.Data.Shop.PaginationPrev.Current}} disabled{{end}}">
						<a class="page-link" href="{{$.Data.Shop.PaginationPrev.Link}}">Previous</a>
					</li>
				{{end}}
				{{range $.Data.Shop.Pagination}}
					{{if .Dots}}
						<li class="page-item disabled"><a class="page-link" href="">...</a></li>
					{{else}}
						<li class="page-item{{if .Current}} active{{end}}">
							<a class="page-link" href="{{.Link}}">{{.Num}}</a>
						</li>
					{{end}}
				{{end}}
				{{if $.Data.Shop.PaginationNext}}
					<li class="page-item{{if $.Data.Shop.PaginationNext.Current}} disabled{{end}}">
						<a class="page-link" href="{{$.Data.Shop.PaginationNext.Link}}">Next</a>
					</li>
				{{end}}
			</ul>
		</nav>
	{{end}}
{{end}}
{{template "footer.html" .}}`)
