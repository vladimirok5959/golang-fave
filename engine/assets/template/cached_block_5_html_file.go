package template

var VarCachedBlock_5HtmlFile = []byte(`{{range $.Data.Blog.Categories 0 0}}
	<li class="{{if and $.Data.Blog.Category (eq $.Data.Blog.Category.Id .Id)}}active{{end}}">
		<a href="{{.Permalink}}">{{.Name}}</a>
	</li>
{{end}}`)
