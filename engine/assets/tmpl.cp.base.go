package assets

var TmplCpBase = []byte(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="theme-color" content="#205081" /><meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"><title>{{$.Data.Title}}</title><link rel="stylesheet" href="{{$.System.PathCssBootstrap}}">{{if or (eq $.System.CpModule "index") (eq $.System.CpModule "blog") (eq $.System.CpModule "shop")}}{{if or (eq $.System.CpSubModule "add") (eq $.System.CpSubModule "modify")}}<link rel="stylesheet" href="{{$.System.PathCssCpWysiwygPell}}">{{end}}{{end}}{{if eq $.System.CpModule "templates"}}<link rel="stylesheet" href="{{$.System.PathCssCpCodeMirror}}">{{end}}<link rel="stylesheet" href="{{$.System.PathCssStyles}}" /><link rel="stylesheet" href="{{$.System.PathCssCpStyles}}"><link rel="shortcut icon" href="{{$.System.PathIcoFav}}" type="image/x-icon" /><script type="text/javascript">var CurrentUserProfileData={first_name:'{{$.Data.UserFirstName}}',last_name:'{{$.Data.UserLastName}}',email:'{{$.Data.UserEmail}}'};function WaitForFave(callback){if(window&&window.fave){callback();}else{setTimeout(function(){WaitForFave(callback);},100);}};</script></head><body class="{{$.Data.BodyClasses}} cp-mod-{{$.System.CpModule}} cp-sub-mod-{{$.System.CpSubModule}}"><div id="sys-modal-user-settings-placeholder"></div><div id="sys-modal-shop-product-attach-placeholder"><div id="sys-modal-files-manager-placeholder"></div></div><div id="sys-modal-system-message-placeholder"></div><nav class="navbar main navbar-expand-md navbar-dark fixed-top bg-dark"><a class="navbar-brand" href="/cp/">{{$.Data.Caption}}</a><button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button><div class="collapse navbar-collapse" id="navbarCollapse"><ul class="navbar-nav mr-auto"><li class="nav-item dropdown"><a class="nav-link dropdown-toggle" href="javascript:;" id="nbModulesDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Modules</a><div class="dropdown-menu" aria-labelledby="nbModulesDropdown">{{$.Data.NavBarModules}}</div></li><li class="nav-item dropdown"><a class="nav-link dropdown-toggle" href="javascript:;" id="nbSystemDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">System</a><div class="dropdown-menu" aria-labelledby="nbSystemDropdown">{{$.Data.NavBarModulesSys}}</div></li><li class="nav-item"><a class="nav-link" href="javascript:fave.FilesManagerDialog();" role="button" aria-expanded="false">Files</a></li></ul><ul class="navbar-nav ml-auto"><li class="nav-item dropdown"><a class="nav-link dropdown-toggle" href="javascript:;" id="nbAccountDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">{{$.Data.UserEmail}}</a><div class="dropdown-menu dropdown-menu-right" aria-labelledby="nbAccountDropdown"><a class="dropdown-item" href="javascript:fave.ModalUserProfile();">My profile</a><div class="dropdown-divider"></div><a class="dropdown-item" href="javascript:fave.ActionLogout('Are you sure want to logout?');">Logout</a></div></li></ul></div></nav><div class="wrap"><div class="sidebar sidebar-left d-block d-sm-table-cell"><div class="scroll"><div class="padd">{{$.Data.SidebarLeft}}</div></div></div><div class="content d-block d-sm-table-cell"><div class="scroll"><div class="padd">{{$.Data.Content}}</div></div></div><div class="sidebar sidebar-right d-none d-lg-table-cell"><div class="scroll"><div class="padd">{{$.Data.SidebarRight}}</div></div></div></div><script src="{{$.System.PathJsJquery}}"></script><script src="{{$.System.PathJsPopper}}"></script><script src="{{$.System.PathJsBootstrap}}"></script>{{if or (eq $.System.CpModule "index") (eq $.System.CpModule "blog") (eq $.System.CpModule "shop")}}{{if or (eq $.System.CpSubModule "add") (eq $.System.CpSubModule "modify")}}<script src="{{$.System.PathJsCpWysiwygPell}}"></script>{{end}}{{end}}{{if eq $.System.CpModule "templates"}}<script src="{{$.System.PathJsCpCodeMirror}}"></script>{{end}}<script src="{{$.System.PathJsCpScripts}}"></script></body></html>`)
