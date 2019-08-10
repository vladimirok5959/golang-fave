package template

var VarSidebarRightHtmlFile = []byte(`<div class="card mb-4">
	<h5 class="card-header">Categories</h5>
	<div class="card-body">
		<ul class="m-0 p-0 pl-4">
			{{range $.Data.Blog.Categories 0}}
				<li class="{{if and $.Data.Blog.Category (eq $.Data.Blog.Category.Id .Id)}}active{{end}}">
					<a href="{{.Permalink}}">{{.Name}}</a>
				</li>
			{{end}}
		</ul>
	</div>
</div>
<div class="card mb-4">
	<h5 class="card-header">Useful links</h5>
	<div class="card-body">
		<ul class="m-0 p-0 pl-4">
			<li><a href="https://github.com/vladimirok5959/golang-fave" target="_blank">Project on GitHub</a></li>
			<li><a href="https://github.com/vladimirok5959/golang-fave/wiki" target="_blank">Wiki on GitHub</a></li>
		</ul>
	</div>
</div>`)