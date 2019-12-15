package ctrlc

import (
	"fmt"
)

func clr(use bool, str string) string {
	if !IS_WIN_PLATFORM && use {
		return fmt.Sprintf("\033[1;31m%s\033[0m", str)
	}
	return str
}

func clg(use bool, str string) string {
	if !IS_WIN_PLATFORM && use {
		return fmt.Sprintf("\033[1;32m%s\033[0m", str)
	}
	return str
}

func cly(use bool, str string) string {
	if !IS_WIN_PLATFORM && use {
		return fmt.Sprintf("\033[1;33m%s\033[0m", str)
	}
	return str
}
