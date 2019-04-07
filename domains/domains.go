package domains

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"golang-fave/utils"
)

type Domains struct {
	sync.RWMutex
	hosts map[string]string
}

func New(www_dir string) *Domains {
	r := Domains{}
	r.hosts = map[string]string{}

	files, err := ioutil.ReadDir(www_dir)
	if err == nil {
		for _, file := range files {
			domains_file := www_dir + string(os.PathSeparator) + file.Name() +
				string(os.PathSeparator) + "config" + string(os.PathSeparator) + ".domains"
			if utils.IsFileExists(domains_file) {
				if f, err := os.Open(domains_file); err == nil {
					defer f.Close()
					reader := bufio.NewReader(f)
					var domain string
					for {
						domain, err = reader.ReadString('\n')
						if err != nil {
							break
						}
						if strings.TrimSpace(domain) != "" {
							r.addDmain(file.Name(), strings.TrimSpace(domain))
						}
					}
				}
			}
		}
	}

	return &r
}

func (this *Domains) addDmain(host string, domain string) {
	this.Lock()
	defer this.Unlock()
	if _, ok := this.hosts[domain]; ok == false {
		this.hosts[domain] = host
	}
}

func (this *Domains) GetHost(domain string) string {
	this.Lock()
	defer this.Unlock()
	if value, ok := this.hosts[domain]; ok == true {
		return value
	}
	return ""
}
