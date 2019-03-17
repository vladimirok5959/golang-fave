package mysqlpool

import (
	"database/sql"
	"sync"
)

type MySqlPool struct {
	sync.RWMutex
	connections map[string]*sql.DB
}

func New() *MySqlPool {
	r := MySqlPool{}
	r.connections = map[string]*sql.DB{}
	return &r
}

func (this *MySqlPool) Get(key string) *sql.DB {
	this.Lock()
	defer this.Unlock()
	if value, ok := this.connections[key]; ok == true {
		return value
	}
	return nil
}

func (this *MySqlPool) Set(key string, value *sql.DB) {
	this.Lock()
	defer this.Unlock()
	this.connections[key] = value
}

func (this *MySqlPool) CloseAll() {
	this.Lock()
	defer this.Unlock()
	for _, c := range this.connections {
		if c != nil {
			c.Close()
		}
	}
}
