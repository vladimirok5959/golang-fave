package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) userToBuffer() {
	if this.bufferUser == nil {
		user := utils.MySql_user{}
		if this.Wrap.S.GetInt("UserId", 0) > 0 {
			err := this.Wrap.DB.QueryRow(`
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
				this.Wrap.S.GetInt("UserId", 0),
			).Scan(
				&user.A_id,
				&user.A_first_name,
				&user.A_last_name,
				&user.A_email,
				&user.A_admin,
				&user.A_active,
			)
			if err != nil {
				this.Wrap.LogError(err.Error())
			}
		}
		this.bufferUser = &user
	}
}

func (this *FERData) UserIsLoggedIn() bool {
	this.userToBuffer()
	return this.bufferUser.A_id > 0
}

func (this *FERData) UserID() int {
	this.userToBuffer()
	return this.bufferUser.A_id
}

func (this *FERData) UserFirstName() string {
	this.userToBuffer()
	return this.bufferUser.A_first_name
}

func (this *FERData) UserLastName() string {
	this.userToBuffer()
	return this.bufferUser.A_last_name
}

func (this *FERData) UserEmail() string {
	this.userToBuffer()
	return this.bufferUser.A_email
}

func (this *FERData) UserIsAdmin() bool {
	this.userToBuffer()
	return this.bufferUser.A_admin > 0
}

func (this *FERData) UserIsActive() bool {
	this.userToBuffer()
	return this.bufferUser.A_active > 0
}
