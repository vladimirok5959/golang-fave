package session

func (this *session) IsSetString(name string) bool {
	if _, ok := this.v.String[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *session) GetString(name string, def string) string {
	if v, ok := this.v.String[name]; ok {
		return v
	} else {
		return def
	}
}

func (this *session) SetString(name string, value string) {
	this.v.String[name] = value
	this.c = true
}
