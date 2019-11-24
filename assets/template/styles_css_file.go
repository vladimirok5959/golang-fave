package template

var VarStylesCssFile = []byte(`/* Fix bootstrap */
@media (max-width: 991.98px) {
	.navbar-expand-lg>.container,
	.navbar-expand-lg>.container-fluid {
		padding-right: 15px;
		padding-left: 15px;
	}
}

@media (max-width: 575px) {
	.navbar-expand-lg>.container,
	.navbar-expand-lg>.container-fluid {
		padding-right: 0px;
		padding-left: 0px;
	}
}

@media (min-width: 992px) {
	.navbar-expand-lg .navbar-nav {
		margin-right: -.5rem;
	}
}

/* Set base bootstrap width */
@media (min-width: 1200px) {
	.container {
		 max-width: 1000px; 
	}

	.navbar-expand-lg .navbar-nav {
		margin-right: -.5rem;
	}
}

/* Base font and colors */
body {
	color: #444;
	font-size: 1.0rem;
	font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
}

.h1,
.h2,
.h3,
.h4,
.h5,
.h6,
h1,
h2,
h3,
h4,
h5,
h6,
footer {
	color: #000;
}

.navbar-brand {
	font-weight: bold;
}

h1 {
	font-size: 250%;
}

/* Nice sticky footer */
html,
body {
	height: 100%;
}

#wrap {
	min-height: 100%;
}

#main {
	overflow: auto;
	padding-bottom: 4.5rem;
}

footer {
	position: relative;
	margin-top: -4.5rem;
	height: 4.5rem;
	clear: both;
}

/* Sticky top nav bar, body class "fixed-top-bar" */
.fixed-top-bar .navbar {
	position: fixed;
	top: 0;
	right: 0;
	left: 0;
	z-index: 1030;
}

.fixed-top-bar #main {
	padding-top: 3.5rem;
}

/* Basket button */
#basket-nav-btn .badge {
	vertical-align: top;
	margin-top: 3px;
}

/* Fave background */
.bg-fave {
	background: #205081 url(/assets/sys/bg.png) repeat 0 0;
}

/* Fix content marging */
.page-content p:last-child {
	margin-bottom: 0px;
}

/* Borders */
.navbar {
	border-bottom: 1px solid rgba(0, 0, 0, .125);
}

footer {
	border-top: 1px solid rgba(0, 0, 0, .125);
}

/* Shop */
.grid-products {
	display: flex;
	flex-flow: row wrap;
	padding-top: 1px;
	padding-left: 1px;
	margin-right: -4px;
}

.grid-products .card-product {
	width: 100%;
	border-radius: 0;
	margin-top: -1px;
	margin-left: -1px;
}

@media (min-width: 768px) {
	.grid-products .card-product {
		width: 50%;
	}
}

@media (min-width: 992px) {
	.grid-products .card-product {
		width: 33.33333%;
	}
}

@media (min-width: 1200px) {
	.grid-products .card-product {
		width: 33.33333%;
	}
}

.grid-products .card-product:hover {
	background-color: #f2f4f6;
}

.grid-products .card-product .card-img-link {
	display: block;
	padding: 1.25rem;
	padding-bottom: 0;
}

.grid-products .card-product .card-img-link img {
	display: block;
	border-radius: 0;
}

.grid-products .card-product .card-title {
	font-size: 1rem;
}

.grid-products .card-product .card-text {
	font-size: 0.7rem;
}

.grid-products .card-product .card-footer {
	padding-top: 0;
	border-radius: 0;
	border-top: 0;
	background: transparent;
	padding-bottom: 1.25rem;
}

.grid-products .card-product .price {
	font-weight: bold;
}

.grid-products .card-product .btn {
	float: right;
}

.product-full .price {
	display: inline-block;
	vertical-align: middle;
}

.product-full .btn-buy {
	display: inline-block;
	vertical-align: middle;
}

.product-full .product-description {
	font-size: 1rem;
}

.product-full .product-description p:last-child {
	margin-bottom: 0px;
}

.product-full .thumbnails {
	padding: 2px;
}

.product-full .thumbnails .thumbnail {
	width: 16.666666667%;
	padding: 2px;
}

.product-full .thumbnails .thumbnail-hidden {
	display: none;
}

.product-full .thumbnails .thumbnail img {
	width: 100%;
	border-radius: 4px;
}

.table-specifications .tcol-1,
.table-specifications .tcol-2 {
	width: 100%;
	display: block;
}

.table-specifications .tcol-1 {
	font-weight: bold;
	border-bottom: none;
	padding-bottom: 0px;
}

.table-specifications .tcol-2 {
	border-top: none;
	padding-top: 0px;
}

@media (min-width: 768px) {
	.table-specifications .tcol-1,
	.table-specifications .tcol-2 {
		width: 50%;
		display: table-cell;
		font-weight: normal;
		border: 1px solid #dee2e6;
		padding: .75rem;
	}
}

@media (min-width: 992px) {
	.table-specifications .tcol-1,
	.table-specifications .tcol-2 {
		width: 50%;
		display: table-cell;
		font-weight: normal;
		border: 1px solid #dee2e6;
		padding: .75rem;
	}
}

@media (min-width: 1200px) {
	.table-specifications .tcol-1,
	.table-specifications .tcol-2 {
		width: 50%;
		display: table-cell;
		font-weight: normal;
		border: 1px solid #dee2e6;
		padding: .75rem;
	}
}

.fixed-top-bar .navbar.navbar-cats {
	position: static;
}

/* Shop basket */
#sys-modal-shop-basket .data .table tbody td {
	vertical-align: middle;
}

#sys-modal-shop-basket .data .table .thc-1 {
	width: 75px;
}

#sys-modal-shop-basket .data .table .thc-4 {
	width: 180px;
	text-align: center;
}

#sys-modal-shop-basket .data .table .thc-4 .btn {
	width: 40px;
}

#sys-modal-shop-basket .data .table .thc-4 .btn-minus {
	float: left;
}

#sys-modal-shop-basket .data .table .thc-4 .btn-plus {
	float: right;
}

#sys-modal-shop-basket .data .table .thc-4 .form-control {
	width: auto;
	display: inline-block;
	text-align: center;
	width: 60px;
}

#sys-modal-shop-basket .data .table .thc-5 {
	width: 110px;
}

#sys-modal-shop-basket .data .table tbody .thc-6 {
	width: 40px;
	font-size: 1.5rem;
	font-weight: bold;
	text-align: right;
	vertical-align: middle;
}

#sys-modal-shop-basket .data .table tbody .thc-6 a:hover {
	text-decoration: none;
}

#sys-modal-shop-basket .data .total {
	text-align: right;
	font-size: 1.5rem;
}

#sys-modal-shop-basket .data .total span {
	display: inline-block;
}

#sys-modal-shop-basket .data .total span.value {
	margin-left: 1rem;
}

#sys-modal-shop-basket .order-form .form-group label {
	font-weight: 700;
	margin-top: .45rem;
	margin-bottom: .45rem;
}

@media (max-width: 768px) {
	#sys-modal-shop-basket .data .table td {
		display: block;
	}

	#sys-modal-shop-basket .data .table .thc-3 {
		display: none;
	}

	#sys-modal-shop-basket .data .table thead {
		display: none;
	}

	#sys-modal-shop-basket .data .table .thc-1 {
		width: auto;
	}

	#sys-modal-shop-basket .data .table .thc-1 img {
		width: 100%;
		height: auto;
	}

	#sys-modal-shop-basket .data .table .thc-4,
	#sys-modal-shop-basket .data .table .thc-5,
	#sys-modal-shop-basket .data .table tbody .thc-6 {
		width: auto;
	}

	#sys-modal-shop-basket .data .table tbody .thc-6 {
		text-align: left;
	}

	#sys-modal-shop-basket .modal-footer {
		display: block;
		text-align: right;
	}

	#sys-modal-shop-basket .modal-footer>:not(:last-child) {
		margin-right: 0px;
		display: block;
	}

	#sys-modal-shop-basket .modal-footer>:not(:first-child) {
		margin-left: 0px;
		margin-top: 1rem;
		display: block;
	}
}

@media (min-width: 768px) {
	#sys-modal-shop-basket .modal-dialog {
		max-width: 660px
	}

	#sys-modal-shop-basket .data .table .thc-3 {
		display: none;
	}
}

@media (min-width: 992px) {
	#sys-modal-shop-basket .modal-dialog {
		max-width: 900px
	}

	#sys-modal-shop-basket .data .table .thc-3 {
		display: table-cell;
	}

	li.currency-changer {
		padding-right: 8px;
	}
}

@media (min-width: 1200px) {
	#sys-modal-shop-basket .modal-dialog {
		max-width: 940px
	}

	#sys-modal-shop-basket .data .table .thc-3 {
		display: table-cell;
	}

	li.currency-changer {
		padding-right: 8px;
	}
}`)
