package resource

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type OneResource struct {
	Path  string
	Ctype string
	Bytes []byte
	MTime int64
}

type Resource struct {
	list map[string]OneResource
}

func New() *Resource {
	r := Resource{}
	r.list = map[string]OneResource{}
	return &r
}

func etag(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func modified(p string, s int, v int64, w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Cache-Control", "no-cache")

	// Set: ETag
	ehash := etag(fmt.Sprintf("%s-%d-%d", p, s, v))
	w.Header().Set("ETag", fmt.Sprintf("%s", ehash))

	// Set: Last-Modified
	w.Header().Set(
		"Last-Modified",
		time.Unix(v, 0).In(time.FixedZone("GMT", 0)).Format("Wed, 01 Oct 2006 15:04:05 GMT"),
	)

	// Check: ETag
	if cc := r.Header.Get("Cache-Control"); cc != "no-cache" {
		if inm := r.Header.Get("If-None-Match"); inm == ehash {
			w.WriteHeader(http.StatusNotModified)
			return false
		}
	}

	// Check: Last-Modified
	if cc := r.Header.Get("Cache-Control"); cc != "no-cache" {
		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if t, err := time.Parse("Wed, 01 Oct 2006 15:04:05 GMT", ims); err == nil {
				if time.Unix(v, 0).In(time.FixedZone("GMT", 0)).Unix() <= t.In(time.FixedZone("GMT", 0)).Unix() {
					w.WriteHeader(http.StatusNotModified)
					return false
				}
			}
		}
	}

	return true
}

func (this *Resource) Add(path string, ctype string, bytes []byte, mtime int64) {
	// Do not add if already in resources list
	if _, ok := this.list[path]; ok == true {
		return
	}

	// Add to resources list
	this.list[path] = OneResource{
		Path:  path,
		Ctype: ctype,
		Bytes: bytes,
		MTime: mtime,
	}
}

func (this *Resource) Response(w http.ResponseWriter, r *http.Request, before func(w http.ResponseWriter, r *http.Request, i *OneResource), after func(w http.ResponseWriter, r *http.Request, i *OneResource)) bool {
	// Do not process if this is not necessary
	if len(r.URL.Path) <= 1 {
		return false
	}

	// Check for resource
	res, ok := this.list[r.URL.Path[1:]]
	if ok == false {
		return false
	}

	// Cache headers
	w.Header().Set("Content-Type", res.Ctype)
	if !modified(r.URL.Path, len(res.Bytes), res.MTime, w, r) {
		return true
	}

	// Call `before` callback
	if before != nil {
		before(w, r, &res)
	}

	// Send resource
	w.Write(res.Bytes)

	// Call `after` callback
	if after != nil {
		after(w, r, &res)
	}

	return true
}
