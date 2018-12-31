package sessions

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Vars struct {
	Int    map[string]int
	String map[string]string
	Bool   map[string]bool
}

type Session struct {
	W            *http.ResponseWriter
	R            *http.Request
	VHost        string
	DirVHostHome string
	RemoteIp     string
	vars         *Vars
	ident        string
	changed      bool
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, vhosthome string, remoteip string) *Session {
	return &Session{w, r, vhost, vhosthome, remoteip, &Vars{}, "", false}
}

func (this *Session) Load() {
	var session, err = this.R.Cookie("fsession")
	if err == nil && len(session.Value) == 40 {
		// Load session
		this.ident = session.Value
		StartNewSession := true
		fsessfile := this.DirVHostHome + "/tmp/" + this.ident
		file, ferr := os.Open(fsessfile)
		if ferr == nil {
			defer file.Close()
			dec := json.NewDecoder(file)
			ferr = dec.Decode(&this.vars)
			if ferr == nil {
				StartNewSession = false
			}
		}
		if StartNewSession {
			sessdata := Vars{}
			sessdata.Int = map[string]int{}
			sessdata.String = map[string]string{}
			sessdata.Bool = map[string]bool{}
			this.vars = &sessdata
			this.changed = true
		}
	} else {
		// Create new session
		// Generate unique hash
		rand.Seed(time.Now().Unix())
		rnd := rand.Intn(9999999-99) + 99
		userstr := this.VHost + this.RemoteIp + this.R.Header.Get("User-Agent") +
			strconv.FormatInt((int64(time.Now().Unix())), 10) +
			strconv.FormatInt(int64(rnd), 10)
		userhashstr := fmt.Sprintf("%x", sha1.Sum([]byte(userstr)))
		this.ident = userhashstr

		// Try to create session file
		sessdata := Vars{}
		sessdata.Int = map[string]int{}
		sessdata.String = map[string]string{}
		sessdata.Bool = map[string]bool{}
		this.vars = &sessdata
		this.changed = true

		// Set session cookie
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "fsession", Value: userhashstr, Expires: expiration}
		http.SetCookie(*this.W, &cookie)
	}
}

func (this *Session) Save() bool {
	if !this.changed {
		return false
	}
	fsessfile := this.DirVHostHome + "/tmp/" + this.ident
	r, err := json.Marshal(this.vars)
	if err == nil {
		file, ferr := os.Create(fsessfile)
		if ferr == nil {
			defer file.Close()
			_, ferr = file.WriteString(string(r))
			if ferr == nil {
				this.changed = false
				return true
			}
		}
	}
	return false
}

func (this *Session) IsSetInt(name string) bool {
	if _, ok := this.vars.Int[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *Session) IsSetString(name string) bool {
	if _, ok := this.vars.String[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *Session) IsSetBool(name string) bool {
	if _, ok := this.vars.Bool[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *Session) SetInt(name string, value int) {
	this.vars.Int[name] = value
	this.changed = true
}

func (this *Session) SetString(name string, value string) {
	this.vars.String[name] = value
	this.changed = true
}

func (this *Session) SetBool(name string, value bool) {
	this.vars.Bool[name] = value
	this.changed = true
}

func (this *Session) GetInt(name string) (int, error) {
	if this.IsSetInt(name) {
		return this.vars.Int[name], nil
	} else {
		return 0, errors.New("Variable is not found")
	}
}

func (this *Session) GetString(name string) (string, error) {
	if this.IsSetString(name) {
		return this.vars.String[name], nil
	} else {
		return "", errors.New("Variable is not found")
	}
}

func (this *Session) GetBool(name string) (bool, error) {
	if this.IsSetBool(name) {
		return this.vars.Bool[name], nil
	} else {
		return false, errors.New("Variable is not found")
	}
}

func (this *Session) GetIntDef(name string, def int) int {
	if this.IsSetInt(name) {
		return this.vars.Int[name]
	} else {
		return def
	}
}

func (this *Session) GetStringDef(name string, def string) string {
	if this.IsSetString(name) {
		return this.vars.String[name]
	} else {
		return def
	}
}

func (this *Session) GetBoolDef(name string, def bool) bool {
	if this.IsSetBool(name) {
		return this.vars.Bool[name]
	} else {
		return def
	}
}

func (this *Session) DelInt(name string) {
	if this.IsSetInt(name) {
		delete(this.vars.Int, name)
		this.changed = true
	}
}

func (this *Session) DelString(name string) {
	if this.IsSetString(name) {
		delete(this.vars.String, name)
		this.changed = true
	}
}

func (this *Session) DelBool(name string) {
	if this.IsSetBool(name) {
		delete(this.vars.Bool, name)
		this.changed = true
	}
}
