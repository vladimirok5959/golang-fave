package builder

import (
	"golang-fave/assets"
)

func CheckBox(value int) string {
	if value == 0 {
		return `<span class="svg-red">` + assets.SysSvgIconError + `</span>`
	} else if value == 1 {
		return `<span class="svg-green">` + assets.SysSvgIconChecked + `</span>`
	} else if value == 2 {
		return `<span class="svg-yellow">` + assets.SysSvgIconRestore + `</span>`
	} else if value == 3 {
		return `<span class="svg-yellow">` + assets.SysSvgIconEmail + `</span>`
	}
	return ""
}
