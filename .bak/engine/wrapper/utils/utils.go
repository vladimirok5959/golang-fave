package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
)

func EmailIsValid(email string) bool {
	regexpe := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regexpe.MatchString(email)
}

func GetMd5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func UrlToArray(url string) []string {
	url_buff := url
	if len(url_buff) >= 1 && url_buff[:1] == "/" {
		url_buff = url_buff[1:]
	}
	if len(url_buff) >= 1 && url_buff[len(url_buff)-1:] == "/" {
		url_buff = url_buff[:len(url_buff)-1]
	}
	if url_buff == "" {
		return []string{}
	} else {
		return strings.Split(url_buff, "/")
	}
}
