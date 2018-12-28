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
	FInt    map[string]int
	FString map[string]string
	FBool   map[string]bool
}

type Session struct {
	w         *http.ResponseWriter
	r         *http.Request
	vhost     string
	vhosthome string
	remoteip  string
	vars      *Vars
	Ident     string
	changed   bool
}

func New(w *http.ResponseWriter, r *http.Request, vhost string, vhosthome string, remoteip string) *Session {
	return &Session{w, r, vhost, vhosthome, remoteip, &Vars{}, "", false}
}

func (s *Session) Load() {
	var session, err = s.r.Cookie("fsession")
	if err == nil && len(session.Value) == 40 {
		// Load session
		s.Ident = session.Value
		StartNewSession := true
		fsessfile := s.vhosthome + "/tmp/" + s.Ident
		file, ferr := os.Open(fsessfile)
		if ferr == nil {
			defer file.Close()
			dec := json.NewDecoder(file)
			ferr = dec.Decode(&s.vars)
			if ferr == nil {
				StartNewSession = false
			}
		}
		if StartNewSession {
			sessdata := Vars{}
			sessdata.FInt = map[string]int{}
			sessdata.FString = map[string]string{}
			sessdata.FBool = map[string]bool{}
			s.vars = &sessdata
			s.changed = true
		}
	} else {
		// Create new session
		// Generate unique hash
		rand.Seed(time.Now().Unix())
		rnd := rand.Intn(9999999-99) + 99
		userstr := s.vhost + s.remoteip + s.r.Header.Get("User-Agent") +
			strconv.FormatInt((int64(time.Now().Unix())), 10) +
			strconv.FormatInt(int64(rnd), 10)
		userhashstr := fmt.Sprintf("%x", sha1.Sum([]byte(userstr)))
		s.Ident = userhashstr

		// Try to create session file
		sessdata := Vars{}
		sessdata.FInt = map[string]int{}
		sessdata.FString = map[string]string{}
		sessdata.FBool = map[string]bool{}
		s.vars = &sessdata
		s.changed = true

		// Set session cookie
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "fsession", Value: userhashstr, Expires: expiration}
		http.SetCookie(*s.w, &cookie)
	}
}

func (s *Session) Save() bool {
	if !s.changed {
		return false
	}
	fsessfile := s.vhosthome + "/tmp/" + s.Ident
	r, err := json.Marshal(s.vars)
	if err == nil {
		file, ferr := os.Create(fsessfile)
		if ferr == nil {
			defer file.Close()
			_, ferr = file.WriteString(string(r))
			if ferr == nil {
				s.changed = false
				return true
			}
		}
	}
	return false
}

func (s *Session) IsSetInt(name string) bool {
	if _, ok := s.vars.FInt[name]; ok {
		return true
	} else {
		return false
	}
}

func (s *Session) IsSetString(name string) bool {
	if _, ok := s.vars.FString[name]; ok {
		return true
	} else {
		return false
	}
}

func (s *Session) IsSetBool(name string) bool {
	if _, ok := s.vars.FBool[name]; ok {
		return true
	} else {
		return false
	}
}

func (s *Session) SetInt(name string, value int) {
	s.vars.FInt[name] = value
	s.changed = true
}

func (s *Session) SetString(name string, value string) {
	s.vars.FString[name] = value
	s.changed = true
}

func (s *Session) SetBool(name string, value bool) {
	s.vars.FBool[name] = value
	s.changed = true
}

func (s *Session) GetInt(name string) (int, error) {
	if s.IsSetInt(name) {
		return s.vars.FInt[name], nil
	} else {
		return 0, errors.New("Variable is not found")
	}
}

func (s *Session) GetString(name string) (string, error) {
	if s.IsSetString(name) {
		return s.vars.FString[name], nil
	} else {
		return "", errors.New("Variable is not found")
	}
}

func (s *Session) GetBool(name string) (bool, error) {
	if s.IsSetBool(name) {
		return s.vars.FBool[name], nil
	} else {
		return false, errors.New("Variable is not found")
	}
}

func (s *Session) DelInt(name string) {
	if s.IsSetInt(name) {
		delete(s.vars.FInt, name)
		s.changed = true
	}
}

func (s *Session) DelString(name string) {
	if s.IsSetString(name) {
		delete(s.vars.FString, name)
		s.changed = true
	}
}

func (s *Session) DelBool(name string) {
	if s.IsSetBool(name) {
		delete(s.vars.FBool, name)
		s.changed = true
	}
}
