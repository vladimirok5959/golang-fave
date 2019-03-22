package mysqlpool

import (
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

func (this *MySqlPool) CloseAll() {
	this.Lock()
	defer this.Unlock()
	for _, c := range this.connections {
		if c != nil {
			c.Close()
		}
	}
}
