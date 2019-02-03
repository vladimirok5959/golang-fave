package session

func (this *session) IsSetBool(name string) bool {
	if _, ok := this.v.Bool[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *session) GetBool(name string, def bool) bool {
	if v, ok := this.v.Bool[name]; ok {
		return v
	} else {
		return def
	}
}

func (this *session) SetBool(name string, value bool) {
	this.v.Bool[name] = value
	this.c = true
}
