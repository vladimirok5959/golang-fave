package builder

import (
	"golang-fave/assets"
)

func CheckBox(value int) string {
	if value > 0 {
		return assets.SysSvgIconChecked
	}
	return assets.SysSvgIconAlert
}
