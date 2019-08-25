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
		/*background: red;*/
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

@media (min-width: 768px) { .grid-products .card-product { width: 50%; } }
@media (min-width: 992px) { .grid-products .card-product { width: 33.33333%; } }
@media (min-width: 1200px) { .grid-products .card-product { width: 33.33333%; } }

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
}`)
