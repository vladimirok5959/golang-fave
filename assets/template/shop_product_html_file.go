package template

var VarShopProductHtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body product-full">
		{{if $.Data.IsUserLoggedIn}}
			{{if $.Data.CurrentUser.IsAdmin}}
				<a href="/cp/shop/modify/{{$.Data.Shop.Product.Id}}/" target="_blank" style="float:right;">Edit</a>
			{{end}}
		{{end}}
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
							<div class="card-body">
								{{if $.Data.Shop.Product.HaveImages}}
									<img class="card-img-top" src="{{$.Data.Shop.Product.Image.Thumbnail3}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
								{{else}}
									<img class="card-img-top" src="{{$.Data.ImagePlaceholderHref}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
								{{end}}
							</div>
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
								<h3>{{if le $.Data.Shop.Product.Quantity 0}}<span class="badge badge-primary">Out of stock</span>{{end}}</h3>
								{{if gt $.Data.Shop.Product.PriceOld 0.00}}<h3 class="price_old mb-0 mr-4"><strike>{{$.Data.Shop.Product.PriceOldNice}} {{$.Data.Shop.CurrentCurrency.Code}}</strike></h3>{{end}}
								<h3 class="price{{if gt $.Data.Shop.Product.PriceOld 0.00}} price_red{{end}} mb-0 mr-4">{{$.Data.Shop.Product.PriceNice}} {{$.Data.Shop.CurrentCurrency.Code}}</h3><button class="btn btn-success btn-buy" onclick="window&&window.frontend&&frontend.ShopBasketProductAdd(this, {{$.Data.Shop.Product.Id}});return false;"{{if le $.Data.Shop.Product.Quantity 0}} disabled{{end}}>Buy</button>
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
										<img class="card-img-top" src="{{$.Data.ImagePlaceholderHref}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
									{{end}}
								</div>
							</div>
							<div class="card mt-3">
								<div class="card-body">
									<h3>{{if le $.Data.Shop.Product.Quantity 0}}<span class="badge badge-primary">Out of stock</span>{{end}}</h3>
									{{if gt $.Data.Shop.Product.PriceOld 0.00}}<h3 class="price_old mb-0 mr-4"><strike>{{$.Data.Shop.Product.PriceOldNice}} {{$.Data.Shop.CurrentCurrency.Code}}</strike></h3>{{end}}
									<h3 class="price{{if gt $.Data.Shop.Product.PriceOld 0.00}} price_red{{end}} mb-0 mr-4">{{$.Data.Shop.Product.PriceNice}} {{$.Data.Shop.CurrentCurrency.Code}}</h3><button class="btn btn-success btn-buy" onclick="window&&window.frontend&&frontend.ShopBasketProductAdd(this, {{$.Data.Shop.Product.Id}});return false;"{{if le $.Data.Shop.Product.Quantity 0}} disabled{{end}}>Buy</button>
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
