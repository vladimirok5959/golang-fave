package wrapper

import (
	"html/template"
)

func (this *Wrapper) ResetCacheBlocks() {
	this.c.Reset(this.Host)
}

func (this *Wrapper) GetBlock1() (template.HTML, bool) {
	return this.c.GetBlock1(this.Host, this.R.URL.Path)
}

func (this *Wrapper) SetBlock1(data template.HTML) {
	this.c.SetBlock1(this.Host, this.R.URL.Path, data)
}

func (this *Wrapper) GetBlock2() (template.HTML, bool) {
	return this.c.GetBlock2(this.Host, this.R.URL.Path)
}

func (this *Wrapper) SetBlock2(data template.HTML) {
	this.c.SetBlock2(this.Host, this.R.URL.Path, data)
}

func (this *Wrapper) GetBlock3() (template.HTML, bool) {
	return this.c.GetBlock3(this.Host, this.R.URL.Path)
}

func (this *Wrapper) SetBlock3(data template.HTML) {
	this.c.SetBlock3(this.Host, this.R.URL.Path, data)
}

func (this *Wrapper) GetBlock4() (template.HTML, bool) {
	return this.c.GetBlock4(this.Host, this.R.URL.Path)
}

func (this *Wrapper) SetBlock4(data template.HTML) {
	this.c.SetBlock4(this.Host, this.R.URL.Path, data)
}

func (this *Wrapper) GetBlock5() (template.HTML, bool) {
	return this.c.GetBlock5(this.Host, this.R.URL.Path)
}

func (this *Wrapper) SetBlock5(data template.HTML) {
	this.c.SetBlock5(this.Host, this.R.URL.Path, data)
}
