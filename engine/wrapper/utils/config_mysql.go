package utils

import (
	"encoding/json"
	"os"
)

type ConfigMySql struct {
	Host     string
	Name     string
	User     string
	Password string
}

func IsMySqlConfigExists(homedir string) bool {
	f, err := os.Open(homedir + "/config/mysql.json")
	if err == nil {
		defer f.Close()
		st, err := os.Stat(homedir + "/config/mysql.json")
		if err == nil {
			if !st.Mode().IsDir() {
				return true
			}
		}
	}
	return false
}

func MySqlConfigRead(homedir string) (*ConfigMySql, error) {
	f, err := os.Open(homedir + "/config/mysql.json")
	if err == nil {
		defer f.Close()
		dec := json.NewDecoder(f)
		conf := ConfigMySql{}
		err = dec.Decode(&conf)
		if err == nil {
			return &conf, err
		}
	}
	return nil, err
}

func MySqlConfigWrite(homedir string, host string, name string, user string, password string) error {
	r, err := json.Marshal(&ConfigMySql{
		Host:     host,
		Name:     name,
		User:     user,
		Password: password,
	})
	if err == nil {
		f, err := os.Create(homedir + "/config/mysql.json")
		if err == nil {
			defer f.Close()
			_, err = f.WriteString(string(r))
			return err
		}
	}
	return err
}
