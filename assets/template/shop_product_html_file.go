package template

var VarShopProductHtmlFile = []byte(`{{template "header.html" .}}
<nav aria-label="breadcrumb">
	<ol class="breadcrumb">
		<li class="breadcrumb-item"><a href="/shop/">Shop</a></li>
		{{if $.Data.Shop.Product.Category.Id}}
			<li class="breadcrumb-item"><a href="{{$.Data.Shop.Product.Category.Permalink}}">{{$.Data.Shop.Product.Category.Name}}</a></li>
		{{end}}
		<li class="breadcrumb-item active"><a href="{{$.Data.Shop.Product.Permalink}}">{{$.Data.Shop.Product.Name}}</a></li>
	</ol>
</nav>
<div class="card mb-4">
	<div class="card-body product-full">
		<h2 class="card-title">{{$.Data.Shop.Product.Name}}</h2>
		<ul class="nav nav-tabs" id="myTab" role="tablist">
			<li class="nav-item">
				<a class="nav-link active" id="all-tab" data-toggle="tab" href="#all" role="tab" aria-controls="all" aria-selected="true">All about product</a>
			</li>
			<li class="nav-item">
				<a class="nav-link" id="specifications-tab" data-toggle="tab" href="#specifications" role="tab" aria-controls="specifications" aria-selected="false">Specifications</a>
			</li>
		</ul>
		<div class="tab-content" id="myTabContent">
			<div class="tab-pane no-fade pt-3 show active" id="all" role="tabpanel" aria-labelledby="all-tab">
				<div class="row">
					<div class="col-md-6">
						<div class="card">
							{{if $.Data.Shop.Product.HaveImages }}
								<img class="card-img-top" src="{{$.Data.Shop.Product.Image.Thumbnail3}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
							{{else}}
								<img class="card-img-top" src="data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22286%22%20height%3D%22180%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20286%20180%22%20preserveAspectRatio%3D%22none%22%3E%3Cdefs%3E%3Cstyle%20type%3D%22text%2Fcss%22%3E%23holder_16c7e5ac360%20text%20%7B%20fill%3Argba(255%2C255%2C255%2C.75)%3Bfont-weight%3Anormal%3Bfont-family%3AHelvetica%2C%20monospace%3Bfont-size%3A14pt%20%7D%20%3C%2Fstyle%3E%3C%2Fdefs%3E%3Cg%20id%3D%22holder_16c7e5ac360%22%3E%3Crect%20width%3D%22286%22%20height%3D%22180%22%20fill%3D%22%23777%22%3E%3C%2Frect%3E%3Cg%3E%3Ctext%20x%3D%22107.0078125%22%20y%3D%2296.234375%22%3E286x180%3C%2Ftext%3E%3C%2Fg%3E%3C%2Fg%3E%3C%2Fsvg%3E" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
							{{end}}
						</div>
					</div>
					<div class="col-md-6">
						<div class="card mt-3 mt-sm-3 mt-md-0 mt-lg-0">
							<div class="card-body">
								<h3 class="price mb-0 mr-4">{{$.Data.Shop.Product.PriceFormat "%.2f"}} {{$.Data.Shop.Product.Currency.Code}}</h3><a href="" class="btn btn-success btn-buy">Buy</a>
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
				
				{{if ne $.Data.Shop.Product.Content ""}}
					<hr>
					<h3>Description</h3>
					<hr>
					<div class="product-description">
						{{$.Data.Shop.Product.Content}}
					</div>
				{{end}}
				<hr>
				<h3>Specifications</h3>
				<hr>
				<!-- <table class="table table-striped table-bordered mb-0 table-specifications">
					<tbody>
						<tr>
							<td class="tcol-1">Диагональ экрана</td>
							<td class="tcol-2">15.6" (1920x1080) Full HD</td>
						</tr>
						<tr>
							<td class="tcol-1">Процессор</td>
							<td class="tcol-2">Шестиядерный Intel Core i7-8750H (2.2 - 4.1 ГГц)</td>
						</tr>
						<tr>
							<td class="tcol-1">Частота обновления экрана</td>
							<td class="tcol-2">60 Гц</td>
						</tr>
						<tr>
							<td class="tcol-1">Объем оперативной памяти</td>
							<td class="tcol-2">8 ГБ</td>
						</tr>
						<tr>
							<td class="tcol-1">Операционная система</td>
							<td class="tcol-2">Windows 10 Home 64bit</td>
						</tr>
						<tr>
							<td class="tcol-1">Цвет</td>
							<td class="tcol-2">Черный</td>
						</tr>
						<tr>
							<td class="tcol-1">Поколение процессора Intel</td>
							<td class="tcol-2">8-ое Coffee Lake</td>
						</tr>
						<tr>
							<td class="tcol-1">Объём накопителя</td>
							<td class="tcol-2">1 ТБ + SSD 128 ГБ</td>
						</tr>
					</tbody>
				</table> -->
			</div>
			<div class="tab-pane no-fade pt-3" id="specifications" role="tabpanel" aria-labelledby="specifications-tab">
				<div class="row">
					<div class="col-md-8">
						<!-- <table class="table table-striped table-bordered mb-0 table-specifications">
							<tbody>
								<tr>
									<td class="tcol-1">Диагональ экрана</td>
									<td class="tcol-2">15.6" (1920x1080) Full HD</td>
								</tr>
								<tr>
									<td class="tcol-1">Процессор</td>
									<td class="tcol-2">Шестиядерный Intel Core i7-8750H (2.2 - 4.1 ГГц)</td>
								</tr>
								<tr>
									<td class="tcol-1">Частота обновления экрана</td>
									<td class="tcol-2">60 Гц</td>
								</tr>
								<tr>
									<td class="tcol-1">Объем оперативной памяти</td>
									<td class="tcol-2">8 ГБ</td>
								</tr>
								<tr>
									<td class="tcol-1">Операционная система</td>
									<td class="tcol-2">Windows 10 Home 64bit</td>
								</tr>
								<tr>
									<td class="tcol-1">Цвет</td>
									<td class="tcol-2">Черный</td>
								</tr>
								<tr>
									<td class="tcol-1">Поколение процессора Intel</td>
									<td class="tcol-2">8-ое Coffee Lake</td>
								</tr>
								<tr>
									<td class="tcol-1">Объём накопителя</td>
									<td class="tcol-2">1 ТБ + SSD 128 ГБ</td>
								</tr>
							</tbody>
						</table> -->
					</div>
					<div class="col-md-4">
						<div class="card mt-3 mt-sm-3 mt-md-0 mt-lg-0">
							<div class="card-body">
								{{if $.Data.Shop.Product.HaveImages }}
									<img class="card-img-top" src="{{$.Data.Shop.Product.Image.Thumbnail2}}" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
								{{else}}
									<img class="card-img-top" src="data:image/svg+xml;charset=UTF-8,%3Csvg%20width%3D%22286%22%20height%3D%22180%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%20286%20180%22%20preserveAspectRatio%3D%22none%22%3E%3Cdefs%3E%3Cstyle%20type%3D%22text%2Fcss%22%3E%23holder_16c7e5ac360%20text%20%7B%20fill%3Argba(255%2C255%2C255%2C.75)%3Bfont-weight%3Anormal%3Bfont-family%3AHelvetica%2C%20monospace%3Bfont-size%3A14pt%20%7D%20%3C%2Fstyle%3E%3C%2Fdefs%3E%3Cg%20id%3D%22holder_16c7e5ac360%22%3E%3Crect%20width%3D%22286%22%20height%3D%22180%22%20fill%3D%22%23777%22%3E%3C%2Frect%3E%3Cg%3E%3Ctext%20x%3D%22107.0078125%22%20y%3D%2296.234375%22%3E286x180%3C%2Ftext%3E%3C%2Fg%3E%3C%2Fg%3E%3C%2Fsvg%3E" alt="{{$.Data.EscapeString $.Data.Shop.Product.Name}}">
								{{end}}
							</div>
						</div>
						<div class="card mt-3">
							<div class="card-body">
								<h3 class="price mb-0 mr-4">{{$.Data.Shop.Product.PriceFormat "%.2f"}} {{$.Data.Shop.Product.Currency.Code}}</h3><a href="" class="btn btn-success btn-buy">Buy</a>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{{template "footer.html" .}}`)
