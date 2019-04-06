package session

func (this *Session) IsSetString(name string) bool {
	if _, ok := this.v.String[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *Session) GetString(name string, def string) string {
	if v, ok := this.v.String[name]; ok {
		return v
	} else {
		return def
	}
}

func (this *Session) SetString(name string, value string) {
	isset := this.IsSetString(name)
	this.v.String[name] = value
	if isset || value != "" {
		this.c = true
	}
}
