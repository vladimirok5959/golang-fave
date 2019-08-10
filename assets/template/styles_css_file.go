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
}`)