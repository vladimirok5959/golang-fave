package session

func (this *session) IsSetInt(name string) bool {
	if _, ok := this.v.Int[name]; ok {
		return true
	} else {
		return false
	}
}

func (this *session) GetInt(name string, def int) int {
	if v, ok := this.v.Int[name]; ok {
		return v
	} else {
		return def
	}
}

func (this *session) SetInt(name string, value int) {
	this.v.Int[name] = value
	this.c = true
}
