package utils

import (
	"encoding/json"
	"os"
)

type ConfigMySql struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func IsMySqlConfigExists(file string) bool {
	return IsFileExists(file) && !IsDir(file)
}

func MySqlConfigRead(file string) (*ConfigMySql, error) {
	f, err := os.Open(file)
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

func MySqlConfigWrite(file string, host string, port string, name string, user string, password string) error {
	r, err := json.Marshal(&ConfigMySql{
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: password,
	})
	if err == nil {
		f, err := os.Create(file)
		if err == nil {
			defer f.Close()
			_, err = f.WriteString(string(r))
			return err
		}
	}
	return err
}
