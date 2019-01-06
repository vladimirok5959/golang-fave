package backend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"

	"golang-fave/constants"
	"golang-fave/engine/backend/modules"
	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

type Backend struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
	user    *utils.MySql_user
	urls    *[]string
}

type TmplData struct {
	Title              string
	BodyClasses        string
	UserId             int
	UserFirstName      string
	UserLastName       string
	UserEmail          string
	UserPassword       string
	UserAvatarLink     string
	ModuleCurrentAlias string
	SidebarLeft        template.HTML
	Content            template.HTML
	SidebarRight       template.HTML
}

func New(wrapper *wrapper.Wrapper, db *sql.DB, url_args *[]string) *Backend {
	return &Backend{wrapper, db, nil, url_args}
}

func (this *Backend) Run() bool {
	// Show add user form if no any user in db
	var count int
	err := this.db.QueryRow("SELECT COUNT(*) FROM `users`;").Scan(&count)
	if this.wrapper.EngineErrMsgOnError(err) {
		return true
	}
	if count <= 0 {
		return this.wrapper.TmplBackEnd(templates.CpFirstUser, nil)
	}

	// Login page
	if this.wrapper.Session.GetIntDef("UserId", 0) <= 0 {
		return this.wrapper.TmplBackEnd(templates.CpLogin, nil)
	}

	// Load current user, if not, show login page
	this.user = &utils.MySql_user{}
	err = this.db.QueryRow("SELECT `id`, `first_name`, `last_name`, `email`, `password` FROM `users` WHERE `id` = ? LIMIT 1;", this.wrapper.Session.GetIntDef("UserId", 0)).Scan(
		&this.user.A_id, &this.user.A_first_name, &this.user.A_last_name, &this.user.A_email, &this.user.A_password)
	if this.wrapper.EngineErrMsgOnError(err) {
		return true
	}
	if this.user.A_id != this.wrapper.Session.GetIntDef("UserId", 0) {
		return this.wrapper.TmplBackEnd(templates.CpLogin, nil)
	}

	// Display cp page
	body_class := "cp"

	// Get module content here
	page_sb_left := ""
	page_content := ""
	page_sb_right := ""

	mdl := modules.New(this.wrapper, this.db, this.user, this.urls)
	if mdl.Run() {
		page_content = mdl.GetContent()
		page_sb_right = mdl.GetSidebarRight()
	}
	page_sb_left = mdl.GetSidebarLeft()

	// If right sidebar and content need to show
	if page_sb_left != "" {
		body_class = body_class + " cp-sidebar-left"
	}
	if page_content == "" {
		body_class = body_class + " cp-404"
		page_content = "Panel 404"
	}
	if page_sb_right != "" {
		body_class = body_class + " cp-sidebar-right"
	}

	// Current module alias
	malias := "index"
	if len(*this.urls) >= 2 {
		malias = (*this.urls)[1]
	}

	// Render page
	page := this.wrapper.TmplParseToString(templates.CpBase, wrapper.TmplDataAll{
		System: this.wrapper.TmplGetSystemData(),
		Data: TmplData{
			Title:              "Fave " + constants.ServerVersion,
			BodyClasses:        body_class,
			UserId:             this.user.A_id,
			UserFirstName:      this.user.A_first_name,
			UserLastName:       this.user.A_last_name,
			UserEmail:          this.user.A_email,
			UserPassword:       "",
			UserAvatarLink:     "https://s.gravatar.com/avatar/" + utils.GetMd5(this.user.A_email) + "?s=80&r=g",
			ModuleCurrentAlias: malias,
			SidebarLeft:        template.HTML(page_sb_left),
			Content:            template.HTML(page_content),
			SidebarRight:       template.HTML(page_sb_right),
		},
	})
	(*this.wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.wrapper.W).Write([]byte(page))
	return true

	return false
}
