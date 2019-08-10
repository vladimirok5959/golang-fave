package template

var VarShopProductHtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body">
		<h2 class="card-title">{{$.Data.Shop.Product.Name}}</h2>
		<div class="page-content">
			{{$.Data.Shop.Product.Briefly}}
			{{$.Data.Shop.Product.Content}}
		</div>
	</div>
	<div class="card-footer text-muted">
		<div>Price: {{$.Data.Shop.Product.PriceFormat "%.2f"}} {{$.Data.Shop.Product.Currency.Code}}</div>
		<div>Published on {{$.Data.Shop.Product.DateTimeFormat "02/01/2006, 15:04:05"}}</div>
		<div>Author: {{$.Data.Shop.Product.User.FirstName}} {{$.Data.Shop.Product.User.LastName}}</div>
	</div>
</div>
{{template "footer.html" .}}`)
