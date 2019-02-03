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
	now := time.Now()
	exp := 7 * 24 * time.Hour
	for _, file := range files {
		if len(file.Name()) == 40 {
			if diff := now.Sub(file.ModTime()); diff > exp {
				err = os.Remove(tmpdir + string(os.PathSeparator) + file.Name())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
