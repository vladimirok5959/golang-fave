package resource

import (
	"net/http"
)

type OneResource struct {
	Path  string
	Ctype string
	Bytes []byte
}

type Resource struct {
	maxurl int
	list   map[string]OneResource
}

func New() *Resource {
	r := Resource{maxurl: 0}
	r.list = map[string]OneResource{}
	return &r
}

func (this *Resource) Add(path string, ctype string, bytes []byte) {
	// Do not add if already in resources list
	if _, ok := this.list[path]; ok == true {
		return
	}

	// Add to resources list
	this.maxurl = len(path)
	this.list[path] = OneResource{
		Path:  path,
		Ctype: ctype,
		Bytes: bytes,
	}
}

func (this *Resource) Response(w http.ResponseWriter, r *http.Request, before func(w http.ResponseWriter, r *http.Request, i *OneResource), after func(w http.ResponseWriter, r *http.Request, i *OneResource)) bool {
	// Do not process if this is not necessary
	if len(r.URL.Path) <= 1 || len(r.URL.Path)-1 > this.maxurl {
		return false
	}

	// Check for resource
	res, ok := this.list[r.URL.Path[1:]]
	if ok == false {
		return false
	}

	// Call `before` callback
	if before != nil {
		before(w, r, &res)
	}

	// Send resource
	w.Header().Set("Content-Type", res.Ctype)
	w.Write(res.Bytes)

	// Call `after` callback
	if after != nil {
		after(w, r, &res)
	}

	return true
}
