package ctrlc

func icon_start(use bool) string {
	if use {
		return C_ICON_START + " "
	}
	return ""
}

func icon_warn(use bool) string {
	if use {
		return C_ICON_WARN + " "
	}
	return ""
}

func icon_hot(use bool) string {
	if use {
		return C_ICON_HOT + " "
	}
	return ""
}

func icon_mag(use bool) string {
	if use {
		return C_ICON_MAG + " "
	}
	return ""
}

func icon_sc(use bool) string {
	if use {
		return C_ICON_SC + " "
	}
	return ""
}
