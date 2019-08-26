package template

var VarCachedBlock_4HtmlFile = []byte(`{{if $.Data.Shop.Category}}
	{{range $.Data.Shop.Categories $.Data.Shop.Category.Id 1}}
		<li class="{{if and $.Data.Shop.Category (eq $.Data.Shop.Category.Id .Id)}}active{{end}}">
			<a href="{{.Permalink}}">{{.Name}}</a>
		</li>
	{{end}}
{{else}}
	{{range $.Data.Shop.Categories 0 1}}
		<li class="{{if and $.Data.Shop.Category (eq $.Data.Shop.Category.Id .Id)}}active{{end}}">
			<a href="{{.Permalink}}">{{.Name}}</a>
		</li>
	{{end}}
{{end}}`)
