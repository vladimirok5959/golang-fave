package wrapper

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func (e *Wrapper) redirectToMainDomain() bool {
	file, err := ioutil.ReadFile(e.DirVhostHome + "/config/domain")
	if err == nil {
		maindomain := strings.TrimSpace(string(file))
		port := ""
		if e.Port != "80" {
			port = ":" + e.Port
		}
		if maindomain+port != e.R.Host {
			http.Redirect(*e.W, e.R, e.R.URL.Scheme+"://"+maindomain+
				port+e.R.URL.RequestURI(), 301)
			return true
		}
	}
	return false
}

func (e *Wrapper) redirectSeoFix() bool {
	full := e.R.URL.RequestURI()
	uris := full[len(e.R.URL.Path):]
	if len(e.R.URL.Path) > 0 {
		if e.R.URL.Path[len(e.R.URL.Path)-1] != '/' {
			http.Redirect(*e.W, e.R, e.R.URL.Path+"/"+uris, 301)
			return true
		}
	} else {
		return false
	}
	return false
}
