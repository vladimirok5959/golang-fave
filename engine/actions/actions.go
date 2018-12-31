package actions

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"golang-fave/engine/wrapper"

	utils "golang-fave/engine/wrapper/utils"
)

type Action struct {
	wrapper *wrapper.Wrapper
	db      *sql.DB
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

func (this *Action) is_valid_email(email string) bool {
	regexpe := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regexpe.MatchString(email)
}

func New(wrapper *wrapper.Wrapper) *Action {
	return &Action{wrapper, nil}
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
