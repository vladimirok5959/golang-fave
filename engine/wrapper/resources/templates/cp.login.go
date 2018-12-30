package templates

var CpLogin = []byte(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="theme-color" content="#205081" /><meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"><title>Please login</title><link rel="stylesheet" href="{{$.System.PathCssBootstrap}}"><link rel="stylesheet" href="{{$.System.PathCssCpStyles}}"><link rel="shortcut icon" href="{{$.System.PathIcoFav}}" type="image/x-icon" /><style type="text/css">html{height:100%;}</style></head><body class="cp-login text-center">{{$.System.BlockModalSysMsg}}<form class="form-signin card" action="/cp/" method="post"><input type="hidden" name="action" value="signin"><h1 class="h3 mb-3 font-weight-normal">Please login</h1><label for="login" class="sr-only">Email address</label><input type="email" id="login" name="login" class="form-control" placeholder="Login" required autofocus><label for="password" class="sr-only">Password</label><input type="password" id="password" name="password" class="form-control mb-3" placeholder="Password" required><button class="btn btn-lg btn-primary btn-block" type="submit">Login</button></form><script src="{{$.System.PathJsJquery}}"></script><script src="{{$.System.PathJsPopper}}"></script><script src="{{$.System.PathJsBootstrap}}"></script><script src="{{$.System.PathJsCpScripts}}"></script></body></html>`)
