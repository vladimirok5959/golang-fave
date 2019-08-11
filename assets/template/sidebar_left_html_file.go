package template

var VarSidebarLeftHtmlFile = []byte(`<div class="card mb-4">
	<h5 class="card-header">Shop categories</h5>
	<div class="card-body">
		<ul class="m-0 p-0 pl-4">
			{{range $.Data.Shop.Categories 0}}
				<li class="{{if and $.Data.Shop.Category (eq $.Data.Shop.Category.Id .Id)}}active{{end}}">
					<a href="{{.Permalink}}">{{.Name}}</a>
				</li>
			{{end}}
		</ul>
	</div>
</div>
<div class="card mb-4">
	<h5 class="card-header">Shop filter</h5>
	<div class="card-body">
		Filter
	</div>
</div>`)
