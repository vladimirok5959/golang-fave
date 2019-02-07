package main

import (
	"io/ioutil"
	"os"
	"time"

	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

func session_clean_start(www_dir string) chan bool {
	ch := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(1 * time.Hour):
				files, err := ioutil.ReadDir(www_dir)
				if err == nil {
					for _, file := range files {
						tmpdir := www_dir + string(os.PathSeparator) + file.Name() + string(os.PathSeparator) + "tmp"
						if utils.IsDirExists(tmpdir) {
							session.Clean(tmpdir)
						}
					}
				}
			case <-ch:
				ch <- true
				return
			}
		}
	}()
	return ch
}

func session_clean_stop(ch chan bool) {
	ch <- true
	<-ch
}
