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
								<img class="card-img-top" src="data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22286%22%20height%3D%22180%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20286%20180%22%20preserveAspectRatio%3D%22none%22%3E%3Cdefs%3E%3Cstyle%20type%3D%22text%2Fcss%22%3E%23holder_16c7e5ac360%20text%20%7B%20fill%3Argba(255%2C255%2C255%2C.75)%3Bfont-weight%3Anormal%3Bfont-family%3AHelvetica%2C%20monospace%3Bfont-size%3A14pt%20%7D%20%3C%2Fstyle%3E%3C%2Fdefs%3E%3Cg%20id%3D%22holder_16c7e5ac360%22%3E%3Crect%20width%3D%22286%22%20height%3D%22180%22%20fill%3D%22%23777%22%3E%3C%2Frect%3E%3Cg%3E%3Ctext%20x%3D%22107.0078125%22%20y%3D%2296.234375%22%3E286x180%3C%2Ftext%3E%3C%2Fg%3E%3C%2Fg%3E%3C%2Fsvg%3E" alt="{{$.Data.EscapeString .Name}}">
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
						<span class="price">{{.PriceFormat "%.2f"}} {{.Currency.Code}}</span><a href="{{.Permalink}}" class="btn btn-success" onclick="window&&window.frontend&&frontend.ShopAddProductToBasket(this, {{.Id}});return false;">Buy</a>
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
