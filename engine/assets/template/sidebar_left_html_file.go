package template

var VarSidebarLeftHtmlFile = []byte(`{{if $.Data.ModuleShopEnabled}}
	<div class="card mb-4">
		<h5 class="card-header">Shop categories</h5>
		<div class="card-body">
			<ul class="m-0 p-0 pl-4">
				{{$.Data.CachedBlock4}}
			</ul>
		</div>
	</div>
	<!-- <div class="card mb-4">
		<h5 class="card-header">Shop filter</h5>
		<div class="card-body">
			Filter
		</div>
	</div> -->
{{end}}`)
