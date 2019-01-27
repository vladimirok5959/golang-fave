package utils

import (
	"os"
	"strings"
)

func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err == nil {
			return true
		}
	}
	return false
}

func IsDir(filename string) bool {
	if st, err := os.Stat(filename); !os.IsNotExist(err) {
		if err == nil {
			if st.Mode().IsDir() {
				return true
			}
		}
	}
	return false
}

func FixPath(path string) string {
	newPath := strings.TrimSpace(path)
	if len(newPath) <= 0 {
		return newPath
	}
	if newPath[len(newPath)-1] == '/' || newPath[len(newPath)-1] == '\\' {
		newPath = newPath[0 : len(newPath)-2]
	}
	return newPath
}
