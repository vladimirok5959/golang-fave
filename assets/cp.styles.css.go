package assets

var CpStylesCss = []byte(`body.cp{padding-right:0!important}.cp .navbar{padding:.5rem 1rem!important}.dropdown-item:focus,.dropdown-item:hover{background-color:#f1f1f1}.dropdown-item.active,.dropdown-item:active{background-color:#007bff}body.cp-login,body.cp-mysql,body.cp-first-user{height:100%;display:-ms-flexbox;display:-webkit-box;display:flex;-ms-flex-align:center;-ms-flex-pack:center;-webkit-box-align:center;align-items:center;-webkit-box-pack:center;justify-content:center;padding-top:40px;padding-bottom:40px;background-color:#f5f5f5}.cp-login .form-signin,.cp-mysql .form-signin,.cp-first-user .form-signin{width:100%;max-width:330px;padding:15px;margin:0 auto}.cp-login .form-signin label,.cp-mysql .form-signin label,.cp-first-user .form-signin label{cursor:pointer}.cp-login .form-signin .form-control,.cp-mysql .form-signin .form-control,.cp-first-user .form-signin .form-control{position:relative;box-sizing:border-box;height:auto;padding:10px;font-size:16px}.cp-login .form-signin .form-control:focus,.cp-mysql .form-signin .form-control:focus,.cp-first-user .form-signin .form-control:focus{z-index:2}.cp-login .form-signin input[type="email"]{margin-bottom:-1px;border-bottom-right-radius:0;border-bottom-left-radius:0}.cp-login .form-signin input[type="password"]{margin-bottom:10px;border-top-left-radius:0;border-top-right-radius:0}.cp-login .sys-messages,.cp-mysql .sys-messages,.cp-first-user .sys-messages{text-align:left}body.cp{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol","Noto Color Emoji";background:initial;background-color:#fff;font-size:1rem;font-weight:400;line-height:1.5;color:#212529;width:100%;height:100%;overflow:hidden}body.cp nav.main{height:56px;box-shadow:0 0 5px 4px rgba(0,0,0,0.25)}body.cp nav.main.bg-dark{background:#0747a6 url(/assets/sys/bg.png) repeat 0 0!important}body.cp nav.main a.navbar-brand{font-weight:700}body.cp nav.main .navbar-nav .nav-item a img{width:35px;height:35px;margin-right:10px;margin-top:-30px;margin-bottom:-30px;background-color:gray}body.cp .wrap{width:100%;height:100%;display:table;align-items:stretch}body.cp .wrap .sidebar,body.cp .wrap .content{display:table-cell;position:relative;padding-top:56px;vertical-align:top}body.cp .wrap .sidebar.sidebar-right{display:none!important}body.cp .wrap .sidebar-right .padd,body.cp .wrap .content .padd{padding:1rem}body.cp .wrap .scroll{height:100%;overflow:hidden;overflow-y:auto}body.cp .wrap .sidebar{width:245px;background:#eee;box-shadow:0 .5em .5em rgba(0,0,0,.3)}body.cp .wrap .sidebar .dropdown-divider{border-color:#d6d6d6;margin:0}body.cp .wrap .sidebar.sidebar-left ul.nav{padding:1rem 0}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item a{color:#444}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item.active{background-color:#417cb9}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item.active a{color:#fff}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item:hover{background-color:#e7e7e7}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item.active:hover{background-color:#417cb9}body.cp .wrap .sidebar.sidebar-left ul.nav ul.nav{background:#eee;padding-top:0}body.cp .wrap .sidebar.sidebar-left ul.nav ul.nav li.nav-item a{color:#444;padding-left:2rem}body.cp .wrap .sidebar.sidebar-left ul.nav ul.nav li.nav-item.active{background-color:#e7e7e7}body.cp .wrap .sidebar.sidebar-left ul.nav ul.nav li.nav-item.active a{color:#417cb9}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item:last-child ul{padding-bottom:0}body.cp .wrap .sidebar.sidebar-left ul.nav li.nav-item svg.sicon{fill:currentColor;margin-right:5px}.svg-green svg{fill:currentColor;color:#28a745}.svg-red svg{fill:currentColor;color:#cb2431}.data-table.table-hover tbody tr:hover{background-color:#fffbdf}.data-table a svg{fill:currentColor;color:#007bff}.data-table a:hover svg{color:#0056b3}.data-table td.col_action a.ico{display:inline-block;width:16px;height:16px;margin-right:10px}.data-table td.col_action a.ico:last-child{margin-right:0}.data-table thead tr{background-color:#e9ecef}.data-table.table-bordered td,.data-table.table-bordered th{border:none;border-top:1px solid #dee2e6}.data-table.table_pages .col_datetime{width:125px}.data-table.table_pages .col_active{width:85px}.data-table.table_pages .col_action{width:100px;text-align:right}.data-table.table_users .col_active,.data-table.table_users .col_admin{width:85px}.data-table.table_users .col_action{width:100px;text-align:right}.data-form label{font-weight:700;margin-top:.45rem;margin-bottom:.45rem}.data-form small{color:#aeb8bc}.data-form > div:nth-last-child(2){margin-bottom:0}.data-form textarea{min-height:5.4rem}.checkbox-ios{display:inline-block}.checkbox-ios input[type=checkbox]{max-height:0;max-width:0;opacity:0;position:absolute}.checkbox-ios input[type=checkbox] + label{display:block;position:relative;box-shadow:inset 0 0 0 1px #ced4da;background:#ced4da;text-indent:-5000px;height:30px;width:60px;border-radius:15px;cursor:pointer;margin-top:0;margin-bottom:0}.checkbox-ios input[type=checkbox] + label:before{content:"";position:absolute;display:block;height:30px;width:30px;top:0;left:0;border-radius:15px;background:rgba(19,191,17,0);-moz-transition:.25s ease-in-out;-webkit-transition:.25s ease-in-out;transition:.25s ease-in-out}.checkbox-ios input[type=checkbox] + label:after{content:"";position:absolute;display:block;height:26px;width:26px;top:2px;left:2px;border-radius:15px;background:#fff;-moz-transition:.25s ease-in-out;-webkit-transition:.25s ease-in-out;transition:.25s ease-in-out}.checkbox-ios input[type=checkbox]:checked + label:before{width:60px;background:#007bff}.checkbox-ios input[type=checkbox]:checked + label:after{left:32px}#sys-modal-user-settings{padding-right:0!important}@media (min-width: 992px){body.cp.cp-sidebar-right .wrap .sidebar.sidebar-right.d-lg-table-cell{display:table-cell!important}}@media (max-width: 575px){body.cp{height:auto;overflow:scroll}body.cp .wrap .sidebar{width:auto;box-shadow:none}body.cp .wrap .content{padding-top:.2rem}}@media (max-width: 767px){.navbar-expand-md .navbar-collapse{padding:1rem;background:#417cb9;box-shadow:0 .2em .2em rgba(0,0,0,.3);border-radius:.25rem}}`)
