package template

var VarEmailNewOrderAdminHtmlFile = []byte(`<html>
	<head>
		<title>{{$.Else.Subject}}</title>
	</head>
	<body>
		<h2>Client</h2>
		<table border="1">
			<tbody>
				<tr>
					<td><b>Last&nbsp;name</b>&nbsp;&nbsp;&nbsp;</td>
					<td>
						{{if ne $.Client.LastName "" }}
							{{$.Client.LastName}}
						{{else}}
							-
						{{end}}
					</td>
				</tr>
				<tr>
					<td><b>First&nbsp;name</b>&nbsp;&nbsp;&nbsp;</td>
					<td>
						{{if ne $.Client.FirstName "" }}
							{{$.Client.FirstName}}
						{{else}}
							-
						{{end}}
					</td>
				</tr>
				<tr>
					<td><b>Middle&nbsp;name</b>&nbsp;&nbsp;&nbsp;</td>
					<td>
						{{if ne $.Client.MiddleName "" }}
							{{$.Client.MiddleName}}
						{{else}}
							-
						{{end}}
					</td>
				</tr>
				<tr>
					<td><b>Phone</b>&nbsp;&nbsp;&nbsp;</td>
					<td>
						{{if ne $.Client.Phone "" }}
							{{$.Client.Phone}}
						{{else}}
							-
						{{end}}
					</td>
				</tr>
				<tr>
					<td><b>Email</b>&nbsp;&nbsp;&nbsp;</td>
					<td>
						{{if ne $.Client.Email "" }}
							{{$.Client.Email}}
						{{else}}
							-
						{{end}}
					</td>
				</tr>
			</tbody>
		</table>
		<div>&nbsp;</div>
		<h2>Delivery</h2>
		<div>
			{{if ne $.Client.DeliveryComment "" }}
				{{$.Client.DeliveryComment}}
			{{else}}
				-
			{{end}}
		</div>
		<div>&nbsp;</div>
		<h2>Order comment</h2>
		<div>
			{{if ne $.Client.OrderComment "" }}
				{{$.Client.OrderComment}}
			{{else}}
				-
			{{end}}
		</div>
		<div>&nbsp;</div>
		<h2>Order products</h2>
		<div>
			<table border="1" width="100%">
				<tbody>
					{{range $.Basket.Products}}
						<tr>
							<td>
								{{.RenderName}}
							</td>
							<td>
								{{.RenderPrice}}&nbsp;{{$.Basket.Currency.Code}}&nbsp;x&nbsp;{{.RenderQuantity}}
							</td>
							<td>
								{{.RenderSum}} {{$.Basket.Currency.Code}}
							</td>
						</tr>
					{{end}}
				</tbody>
			</table>
		</div>
		<h2>Total: {{$.Basket.RenderTotalSum}} {{$.Basket.Currency.Code}}</h2>
		<div>&nbsp;</div>
		<div><a href="{{$.Else.CpOrderLink}}" target="_blank">{{$.Else.CpOrderLink}}</a></div>
	</body>
</html>`)
