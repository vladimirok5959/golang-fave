package session

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type vars struct {
	Bool   map[string]bool
	Int    map[string]int
	String map[string]string
}

type Session struct {
	w http.ResponseWriter
	r *http.Request
	d string
	v *vars
	c bool
	i string
}

func New(w http.ResponseWriter, r *http.Request, tmpdir string) *Session {
	sess := Session{w: w, r: r, d: tmpdir, v: &vars{}, c: false, i: ""}

	cookie, err := r.Cookie("session")
	if err == nil && len(cookie.Value) == 40 {
		// Load from file
		sess.i = cookie.Value
		fname := sess.d + string(os.PathSeparator) + sess.i
		f, err := os.Open(fname)
		if err == nil {
			defer f.Close()
			dec := json.NewDecoder(f)
			err = dec.Decode(&sess.v)
			if err == nil {
				return &sess
			}

			// Update file last modify time if needs
			if info, err := os.Stat(fname); err == nil {
				if time.Now().Sub(info.ModTime()) > 30*time.Minute {
					_ = os.Chtimes(fname, time.Now(), time.Now())
				}
			}
		}
	} else {
		// Create new
		rand.Seed(time.Now().Unix())

		// Real remote IP for proxy servers
		rRemoteAddr := r.RemoteAddr
		if r.Header.Get("X-Real-IP") != "" && len(r.Header.Get("X-Real-IP")) <= 25 {
			rRemoteAddr = rRemoteAddr + ", " + strings.TrimSpace(r.Header.Get("X-Real-IP"))
		} else if r.Header.Get("X-Forwarded-For") != "" && len(r.Header.Get("X-Forwarded-For")) <= 25 {
			rRemoteAddr = rRemoteAddr + ", " + strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
		}

		sign := rRemoteAddr + r.Header.Get("User-Agent") + fmt.Sprintf("%d", int64(time.Now().Unix())) + fmt.Sprintf("%d", int64(rand.Intn(9999999-99)+99))
		sess.i = fmt.Sprintf("%x", sha1.Sum([]byte(sign)))

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sess.i,
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
		})
	}

	// Init empty
	sess.v = &vars{
		Bool:   map[string]bool{},
		Int:    map[string]int{},
		String: map[string]string{},
	}

	return &sess
}

func (this *Session) Close() bool {
	if !this.c {
		return false
	}

	r, err := json.Marshal(this.v)
	if err == nil {
		f, err := os.Create(this.d + string(os.PathSeparator) + this.i)
		if err == nil {
			defer f.Close()
			_, err = f.Write(r)
			if err == nil {
				this.c = false
				return true
			}
		}
	}

	return false
}

func (this *Session) Destroy() error {
	if this.d == "" || this.i == "" {
		return nil
	}
	err := os.Remove(this.d + string(os.PathSeparator) + this.i)
	if err == nil {
		this.c = false
	}
	return err
}
