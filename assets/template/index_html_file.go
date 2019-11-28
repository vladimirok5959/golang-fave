package template

var VarIndexHtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body">
		{{if $.Data.IsUserLoggedIn}}
			{{if $.Data.CurrentUser.IsAdmin}}
				<a href="/cp/index/modify/{{$.Data.Page.Id}}/" target="_blank" style="float:right;">Edit</a>
			{{end}}
		{{end}}
		<h2 class="card-title">{{$.Data.Page.Name}}</h2>
		<div class="page-content">
			{{$.Data.Page.Content}}
		</div>
	</div>
	<div class="card-footer text-muted">
		<div>Published on {{$.Data.Page.DateTimeFormat "02/01/2006, 15:04:05"}}</div>
		<div>Author: {{$.Data.Page.User.FirstName}} {{$.Data.Page.User.LastName}}</div>
	</div>
</div>
{{template "footer.html" .}}`)
