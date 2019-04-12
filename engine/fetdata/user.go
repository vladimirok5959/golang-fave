package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type User struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_user
}

func (this *User) load(id int) {
	if this == nil {
		return
	}
	if this.object != nil {
		return
	}
	this.object = &utils.MySql_user{}
	if err := this.wrap.DB.QueryRow(`
		SELECT
			id,
			first_name,
			last_name,
			email,
			admin,
			active
		FROM
			users
		WHERE
			id = ?
		LIMIT 1;`,
		id,
	).Scan(
		&this.object.A_id,
		&this.object.A_first_name,
		&this.object.A_last_name,
		&this.object.A_email,
		&this.object.A_admin,
		&this.object.A_active,
	); err != nil {
		return
	}
}

func (this *User) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *User) FirstName() string {
	if this == nil {
		return ""
	}
	return this.object.A_first_name
}

func (this *User) LastName() string {
	if this == nil {
		return ""
	}
	return this.object.A_last_name
}

func (this *User) Email() string {
	if this == nil {
		return ""
	}
	return this.object.A_email
}

func (this *User) IsAdmin() bool {
	if this == nil {
		return false
	}
	return this.object.A_admin == 1
}

func (this *User) IsActive() bool {
	if this == nil {
		return false
	}
	return this.object.A_active == 1
}
