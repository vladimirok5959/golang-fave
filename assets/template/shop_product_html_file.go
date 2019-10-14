package template

var VarShopProductHtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body product-full">
		<h2 class="card-title">{{$.Data.Shop.Product.Name}} {{$.Data.Shop.Product.Id}}</h2>
		<ul class="nav nav-tabs" id="myTab" role="tablist">
			<li class="nav-item">
				<a class="nav-link active" id="all-tab" data-toggle="tab" href="#all" role="tab" aria-controls="all" aria-selected="true">All about product</a>
			</li>
			{{if $.Data.Shop.Product.HaveSpecs}}
				<li class="nav-item">
					<a class="nav-link" id="specifications-tab" data-toggle="tab" href="#specifications" role="tab" aria-controls="specifications" aria-selected="false">Specifications</a>
				</li>
			{{end}}
		</ul>
		<div class="tab-content" id="myTabContent">
			<div class="tab-pane no-fade pt-3 show active" id="all" role="tabpanel" aria-labelledby="all-tab">
				<div class="row">
					<div class="col-md-6">
						<div class="card" id="product_image">
							{{if $.Data.Shop.Product.HaveImages}}
								<img class="card-img-top" src="{{$.Data.Shop.Product.Image.Thumbnail3}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
							{{else}}
								<img class="card-img-top" src="data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22286%22%20height%3D%22180%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20286%20180%22%20preserveAspectRatio%3D%22none%22%3E%3Cdefs%3E%3Cstyle%20type%3D%22text%2Fcss%22%3E%23holder_16c7e5ac360%20text%20%7B%20fill%3Argba(255%2C255%2C255%2C.75)%3Bfont-weight%3Anormal%3Bfont-family%3AHelvetica%2C%20monospace%3Bfont-size%3A14pt%20%7D%20%3C%2Fstyle%3E%3C%2Fdefs%3E%3Cg%20id%3D%22holder_16c7e5ac360%22%3E%3Crect%20width%3D%22286%22%20height%3D%22180%22%20fill%3D%22%23777%22%3E%3C%2Frect%3E%3Cg%3E%3Ctext%20x%3D%22107.0078125%22%20y%3D%2296.234375%22%3E286x180%3C%2Ftext%3E%3C%2Fg%3E%3C%2Fg%3E%3C%2Fsvg%3E" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
							{{end}}
						</div>
						{{if $.Data.Shop.Product.HaveImages}}
							<div class="card mt-1">
								<div id="product_thumbnails" class="thumbnails d-flex flex-wrap">
									{{range $index, $img := $.Data.Shop.Product.Images}}
										<a class="thumbnail{{if gt $index 5}} thumbnail-hidden{{end}}" href="{{.ThumbnailFull}}" data-src="{{.ThumbnailFull}}" data-hover="{{.Thumbnail3}}" data-index="{{$index}}">
											<img class="img-responsive" alt="" src="{{.Thumbnail0}}" />
										</a>
									{{end}}
								</div>
							</div>
						{{end}}
					</div>
					<div class="col-md-6">
						{{if $.Data.Shop.Product.HaveVariations}}
							<div class="card mt-3 mt-sm-3 mt-md-0 mt-lg-0">
								<select class="form-control" onchange="document.location=this.value;">
									{{range $variation := $.Data.Shop.Product.Variations}}
										<option value="{{.Link}}"{{if .Selected}} selected{{end}}>{{.Name}}</option>
									{{end}}
								</select>
							</div>
						{{end}}
						<div class="card mt-3{{if not $.Data.Shop.Product.HaveVariations}} mt-sm-3 mt-md-0 mt-lg-0{{end}}">
							<div class="card-body">
								<h3 class="price mb-0 mr-4">{{$.Data.Shop.Product.PriceFormat "%.2f"}} {{$.Data.Shop.Product.Currency.Code}}</h3><a href="" class="btn btn-success btn-buy" onclick="window&&window.frontend&&frontend.ShopBasketProductPlus(this, {{$.Data.Shop.Product.Id}});return false;">Buy</a>
							</div>
						</div>
						<div class="card mt-3">
							<div class="card-header">Payment</div>
							<div class="card-body">
								<p class="card-text">Non-cash, Cash</p>
							</div>
						</div>
					</div>
				</div>
				<div class="row">
					<div class="col-md-12">
						{{if ne $.Data.Shop.Product.Content ""}}
							<hr>
							<h3>Description</h3>
							<hr>
							<div class="product-description">
								{{$.Data.Shop.Product.Content}}
							</div>
						{{end}}
						{{if $.Data.Shop.Product.HaveSpecs}}
							<hr>
							<h3>Specifications</h3>
							<hr>
							<table class="table table-striped table-bordered mb-0 table-specifications">
								<tbody>
									{{range $.Data.Shop.Product.Specs}}
										<tr>
											<td class="tcol-1">{{.FilterName}}</td>
											<td class="tcol-2">{{.FilterValue}}</td>
										</tr>
									{{end}}
								</tbody>
							</table>
						{{end}}
					</div>
				</div>
			</div>
			{{if $.Data.Shop.Product.HaveSpecs}}
				<div class="tab-pane no-fade pt-3" id="specifications" role="tabpanel" aria-labelledby="specifications-tab">
					<div class="row">
						<div class="col-md-8">
							{{if $.Data.Shop.Product.HaveSpecs}}
								<table class="table table-striped table-bordered mb-0 table-specifications">
									<tbody>
										{{range $.Data.Shop.Product.Specs}}
											<tr>
												<td class="tcol-1">{{.FilterName}}</td>
												<td class="tcol-2">{{.FilterValue}}</td>
											</tr>
										{{end}}
									</tbody>
								</table>
							{{end}}
						</div>
						<div class="col-md-4">
							<div class="card mt-3 mt-sm-3 mt-md-0 mt-lg-0">
								<div class="card-body">
									{{if $.Data.Shop.Product.HaveImages}}
										<img class="card-img-top" src="{{$.Data.Shop.Product.Image.Thumbnail2}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
									{{else}}
										<img class="card-img-top" src="data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22286%22%20height%3D%22180%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20286%20180%22%20preserveAspectRatio%3D%22none%22%3E%3Cdefs%3E%3Cstyle%20type%3D%22text%2Fcss%22%3E%23holder_16c7e5ac360%20text%20%7B%20fill%3Argba(255%2C255%2C255%2C.75)%3Bfont-weight%3Anormal%3Bfont-family%3AHelvetica%2C%20monospace%3Bfont-size%3A14pt%20%7D%20%3C%2Fstyle%3E%3C%2Fdefs%3E%3Cg%20id%3D%22holder_16c7e5ac360%22%3E%3Crect%20width%3D%22286%22%20height%3D%22180%22%20fill%3D%22%23777%22%3E%3C%2Frect%3E%3Cg%3E%3Ctext%20x%3D%22107.0078125%22%20y%3D%2296.234375%22%3E286x180%3C%2Ftext%3E%3C%2Fg%3E%3C%2Fg%3E%3C%2Fsvg%3E" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
									{{end}}
								</div>
							</div>
							<div class="card mt-3">
								<div class="card-body">
									<h3 class="price mb-0 mr-4">{{$.Data.Shop.Product.PriceFormat "%.2f"}} {{$.Data.Shop.Product.Currency.Code}}</h3><a href="" class="btn btn-success btn-buy" onclick="window&&window.frontend&&frontend.ShopBasketProductPlus(this, {{$.Data.Shop.Product.Id}});return false;">Buy</a>
								</div>
							</div>
						</div>
					</div>
				</div>
			{{end}}
		</div>
	</div>
</div>
{{template "footer.html" .}}`)
