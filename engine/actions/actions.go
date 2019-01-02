package actions

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
)

type Action struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
	user    *utils.MySql_user
}

func (this *Action) write(data string) {
	(*this.wrapper.W).Write([]byte(data))
}

func (this *Action) msg_show(title string, msg string) {
	this.write(fmt.Sprintf(
		`ModalShowMsg('%s', '%s');`,
		strings.Replace(strings.Replace(title, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1),
		strings.Replace(strings.Replace(msg, `'`, `&rsquo;`, -1), `"`, `&rdquo;`, -1)))
}

func (this *Action) msg_success(msg string) {
	this.msg_show("Success", msg)
}

func (this *Action) msg_error(msg string) {
	this.msg_show("Error", msg)
}

func (this *Action) use_database() error {
	if this.db != nil {
		return errors.New("already connected to database")
	}
	if !utils.IsMySqlConfigExists(this.wrapper.DirVHostHome) {
		return errors.New("can't read database configuration file")
	}
	mc, err := utils.MySqlConfigRead(this.wrapper.DirVHostHome)
	if err != nil {
		return err
	}
	this.db, err = sql.Open("mysql", mc.User+":"+mc.Password+"@tcp("+mc.Host+":"+mc.Port+")/"+mc.Name)
	if err != nil {
		return err
	}
	err = this.db.Ping()
	if err != nil {
		this.db.Close()
		return err
	}
	return nil
}

func (this *Action) load_session_user() error {
	if this.db == nil {
		return errors.New("not connected to database")
	}
	if this.user != nil {
		return errors.New("user already loaded")
	}
	if this.wrapper.Session.GetIntDef("UserId", 0) <= 0 {
		return errors.New("session user id is not defined")
	}
	this.user = &utils.MySql_user{}
	err := this.db.QueryRow("SELECT `id`, `first_name`, `last_name`, `email`, `password` FROM `users` WHERE `id` = ? LIMIT 1;", this.wrapper.Session.GetIntDef("UserId", 0)).Scan(
		&this.user.A_id, &this.user.A_first_name, &this.user.A_last_name, &this.user.A_email, &this.user.A_password)
	if err != nil {
		return err
	}
	if this.user.A_id != this.wrapper.Session.GetIntDef("UserId", 0) {
		return errors.New("can't load user from session user id")
	}
	return nil
}

func New(wrapper *wrapper.Wrapper) *Action {
	return &Action{wrapper, nil, nil}
}

func (this *Action) Run() bool {
	if this.wrapper.R.Method != "POST" {
		return false
	}
	if err := this.wrapper.R.ParseForm(); err == nil {
		action := this.wrapper.R.FormValue("action")
		if action != "" {
			if _, ok := reflect.TypeOf(this).MethodByName("Action_" + action); ok {
				(*this.wrapper.W).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				(*this.wrapper.W).Header().Set("Content-Type", "text/html; charset=utf-8")
				reflect.ValueOf(this).MethodByName("Action_" + action).Call([]reflect.Value{})
				return true
			}
		}
	}
	return false
}
