package wrapper

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func (this *Wrapper) redirectToMainDomain() bool {
	file, err := ioutil.ReadFile(this.DirVHostHome + "/config/domain")
	if err == nil {
		maindomain := strings.TrimSpace(string(file))
		port := ""
		if this.Port != "80" {
			port = ":" + this.Port
		}
		if maindomain+port != this.R.Host {
			http.Redirect(*this.W, this.R, this.R.URL.Scheme+"://"+maindomain+
				port+this.R.URL.RequestURI(), 301)
			return true
		}
	}
	return false
}

func (this *Wrapper) redirectSeoFix() bool {
	full := this.R.URL.RequestURI()
	uris := full[len(this.R.URL.Path):]
	if len(this.R.URL.Path) > 0 {
		if this.R.URL.Path[len(this.R.URL.Path)-1] != '/' {
			http.Redirect(*this.W, this.R, this.R.URL.Path+"/"+uris, 301)
			return true
		}
	} else {
		return false
	}
	return false
}
