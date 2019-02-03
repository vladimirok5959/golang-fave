package static

import (
	"net/http"
	"os"
)

type Static struct {
	dirIndexFile string
}

func New(file string) *Static {
	r := Static{dirIndexFile: file}
	return &r
}

func (this *Static) Response(dir string, w http.ResponseWriter, r *http.Request, before func(w http.ResponseWriter, r *http.Request), after func(w http.ResponseWriter, r *http.Request)) bool {
	if r.URL.Path == "/" {
		f, err := os.Open(dir + string(os.PathSeparator) + this.dirIndexFile)
		if err == nil {
			defer f.Close()
			st, err := os.Stat(dir + string(os.PathSeparator) + this.dirIndexFile)
			if err != nil {
				return false
			}
			if st.Mode().IsDir() {
				return false
			}
			if before != nil {
				before(w, r)
			}
			http.ServeFile(w, r, dir+string(os.PathSeparator)+this.dirIndexFile)
			if after != nil {
				after(w, r)
			}
			return true
		}
	} else {
		f, err := os.Open(dir + r.URL.Path)
		if err == nil {
			defer f.Close()
			st, err := os.Stat(dir + r.URL.Path)
			if err != nil {
				return false
			}
			if st.Mode().IsDir() {
				if r.URL.Path[len(r.URL.Path)-1] == '/' {
					fi, err := os.Open(dir + r.URL.Path + string(os.PathSeparator) + this.dirIndexFile)
					if err == nil {
						defer fi.Close()
						sti, err := os.Stat(dir + r.URL.Path + string(os.PathSeparator) + this.dirIndexFile)
						if err != nil {
							return false
						}
						if sti.Mode().IsDir() {
							return false
						}
						if before != nil {
							before(w, r)
						}
						http.ServeFile(w, r, dir+r.URL.Path+string(os.PathSeparator)+this.dirIndexFile)
						if after != nil {
							after(w, r)
						}
						return true
					}
				}
				return false
			}
			if before != nil {
				before(w, r)
			}
			http.ServeFile(w, r, dir+r.URL.Path)
			if after != nil {
				after(w, r)
			}
			return true
		}
	}
	return false
}
