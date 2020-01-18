package mysqlpool

import (
	"errors"
	"strings"
	"sync"

	"golang-fave/engine/sqlw"
)

type MySqlPool struct {
	sync.RWMutex
	connections map[string]*sqlw.DB
}

func New() *MySqlPool {
	r := MySqlPool{}
	r.connections = map[string]*sqlw.DB{}
	return &r
}

func (this *MySqlPool) Get(key string) *sqlw.DB {
	this.Lock()
	defer this.Unlock()
	if value, ok := this.connections[key]; ok == true {
		return value
	}
	return nil
}

func (this *MySqlPool) Set(key string, value *sqlw.DB) {
	this.Lock()
	defer this.Unlock()
	this.connections[key] = value
}

func (this *MySqlPool) Del(key string) {
	if _, ok := this.connections[key]; ok {
		delete(this.connections, key)
	}
}

func (this *MySqlPool) Close() error {
	this.Lock()
	defer this.Unlock()
	var errs []string
	for _, c := range this.connections {
		if c != nil {
			if err := c.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}
