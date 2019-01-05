package backend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"

	"golang-fave/constants"
	"golang-fave/engine/wrapper"

	templates "golang-fave/engine/wrapper/resources/templates"
	utils "golang-fave/engine/wrapper/utils"
)

type Backend struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
	user    *utils.MySql_user
}

type TmplData struct {
	Title          string
	UserId         int
	UserFirstName  string
	UserLastName   string
	UserEmail      string
	UserPassword   string
	UserAvatarLink string
	SidebarLeft    template.HTML
	Content        template.HTML
	SidebarRight   template.HTML
}

func New(wrapper *wrapper.Wrapper, db *sql.DB) *Backend {
	return &Backend{wrapper, db, nil}
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
	/*
		(*this.wrapper.W).Write([]byte(`Admin panel here...`))
		return true
	*/
	// return this.wrapper.TmplBackEnd(templates.CpBase, nil)

	/*
		tmpl, err := template.New("template").Parse(string(templates.CpBase))
		if err == nil {
			(*this.wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			tmpl.Execute(*this.wrapper.W, wrapper.TmplDataAll{
				System: this.wrapper.TmplGetSystemData(),
				Data: TmplData{
					Title: "Fave " + constants.ServerVersion,
					UserEmail: this.user.A_email,
					SidebarLeft: "Sidebar left",
					Content: "Content",
					SidebarRight: "Sidebar right",
				},
			})
			return true
		}
	*/

	// Get parsed template as string
	// https://coderwall.com/p/ns60fq/simply-output-go-html-template-execution-to-strings

	// http://localhost:8080/admin/

	sidebar_left := string(`<ul class="nav flex-column">
		<li class="nav-item active">
			<a class="nav-link" href="#">Pages</a>
			<ul class="nav flex-column">
				<li class="nav-item active">
					<a class="nav-link" href="#">List of pages</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="#">Add new page</a>
				</li>
			</ul>
		</li>
		<li class="nav-item">
			<a class="nav-link" href="#">Link 2</a>
		</li>
		<li class="nav-item">
			<a class="nav-link" href="#">Link 3</a>
		</li>
		<li class="nav-item">
			<a class="nav-link" href="#">Link 4</a>
		</li>
	</ul>`)

	page := this.wrapper.TmplParseToString(templates.CpBase, wrapper.TmplDataAll{
		System: this.wrapper.TmplGetSystemData(),
		Data: TmplData{
			Title:          "Fave " + constants.ServerVersion,
			UserId:         this.user.A_id,
			UserFirstName:  this.user.A_first_name,
			UserLastName:   this.user.A_last_name,
			UserEmail:      this.user.A_email,
			UserPassword:   "",
			UserAvatarLink: "https://s.gravatar.com/avatar/" + utils.GetMd5(this.user.A_email) + "?s=80&r=g",
			SidebarLeft:    template.HTML(sidebar_left),
			Content:        template.HTML("Content"),
			SidebarRight:   template.HTML("Sidebar right"),
		},
	})
	(*this.wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*this.wrapper.W).Write([]byte(page))
	return true

	return false
}
