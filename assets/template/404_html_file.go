package template

var Var404HtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body">
		<h2 class="card-title">Error 404</h2>
		<div class="page-content">
			The page what you looking for "<b>{{$.Data.RequestURL}}</b>" is not found
		</div>
	</div>
</div>
{{template "footer.html" .}}`)