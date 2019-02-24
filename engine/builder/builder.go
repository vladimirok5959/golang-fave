package builder

import (
	"golang-fave/assets"
)

func CheckBox(value int) string {
	if value > 0 {
		return `<span class="svg-green">` + assets.SysSvgIconChecked + `</span>`
	}
	return `<span class="svg-red">` + assets.SysSvgIconError + `</span>`
}
