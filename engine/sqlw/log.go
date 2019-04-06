package sqlw

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang-fave/consts"
)

func log(query string, s time.Time, e error, transaction bool) {
	msg := query
	if reg, err := regexp.Compile("[\\s\\t]+"); err == nil {
		msg = strings.Trim(reg.ReplaceAllString(msg, " "), " ")
	}
	if reg, err := regexp.Compile("[\\s\\t]+;$"); err == nil {
		msg = reg.ReplaceAllString(msg, ";")
	}
	eStr := " (nil)"
	if e != nil {
		eStr = " (" + e.Error() + ") "
	}
	if consts.IS_WIN {
		fmt.Fprintln(os.Stdout, "[SQL] "+msg+eStr+fmt.Sprintf(" %.3f ms", time.Now().Sub(s).Seconds()))
	} else {
		color := "0;33"
		if transaction {
			color = "1;33"
		}
		fmt.Fprintln(os.Stdout, "\033["+color+"m[SQL] "+msg+eStr+fmt.Sprintf(" %.3f ms", time.Now().Sub(s).Seconds())+"\033[0m")
	}
}
