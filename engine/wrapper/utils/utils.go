package utils

import (
	"regexp"
)

func EmailIsValid(email string) bool {
	regexpe := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regexpe.MatchString(email)
}
