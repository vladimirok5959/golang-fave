package template

var VarFooterHtmlFile = []byte(`						</div>
						<div class="col-md-4">
							{{template "sidebar-right.html" .}}
						</div>
					</div>
				</div>
			</div>
		</div>
		<footer class="bg-light py-4">
			<div class="container">
				<p class="m-0 text-center text-black">
					Copyright © Your Website {{if eq ($.Data.DateTimeFormat "2006") "2019"}}
						{{$.Data.DateTimeFormat "2006"}}
					{{else}}
						2019-{{$.Data.DateTimeFormat "2006"}}
					{{end}}
				</p>
			</div>
		</footer>
		<!-- Optional JavaScript -->
		<!-- jQuery first, then Popper.js, then Bootstrap JS -->
		<script src="{{$.System.PathJsJquery}}"></script>
		<script src="{{$.System.PathJsPopper}}"></script>
		<script src="{{$.System.PathJsBootstrap}}"></script>
		<script src="{{$.System.PathThemeScripts}}"></script>
	</body>
</html>`)
