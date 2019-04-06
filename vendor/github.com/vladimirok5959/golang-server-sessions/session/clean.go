package session

import (
	"io/ioutil"
	"os"
	"time"
)

func Clean(tmpdir string) error {
	files, err := ioutil.ReadDir(tmpdir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if len(file.Name()) == 40 {
			if diff := time.Now().Sub(file.ModTime()); diff > 24*time.Hour {
				err = os.Remove(tmpdir + string(os.PathSeparator) + file.Name())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
