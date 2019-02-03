package session

func (this *Session) IsSetBool(name string) bool {
	if _, ok := this.v.Bool[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *Session) GetBool(name string, def bool) bool {
	if v, ok := this.v.Bool[name]; ok {
		return v
	} else {
		return def
	}
}

func (this *Session) SetBool(name string, value bool) {
	this.v.Bool[name] = value
	this.c = true
}
