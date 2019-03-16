package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang-fave/utils"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

func session_clean_start(www_dir string) (chan bool, chan bool) {
	ch := make(chan bool)
	stop := make(chan bool)
	go func() {
		for {
			// Destroy old session files on each host
			// Every one hour
			select {
			case <-time.After(1 * time.Hour):
				files, err := ioutil.ReadDir(www_dir)
				if err == nil {
					for _, file := range files {
						select {
						case <-stop:
							break
						default:
							tmpdir := www_dir + string(os.PathSeparator) + file.Name() + string(os.PathSeparator) + "tmp"
							if utils.IsDirExists(tmpdir) {
								session.Clean(tmpdir)
							}
						}
					}
				}
			case <-ch:
				ch <- true
				return
			}
		}
	}()
	return ch, stop
}

func session_clean_stop(ch, stop chan bool) {
	for {
		select {
		case stop <- true:
		case ch <- true:
			<-ch
			return
		case <-time.After(8 * time.Second):
			fmt.Println("Session cleaner error: force exit by timeout after 8 seconds")
			return
		}
	}
}
