package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) ResetBlock4(host string) {
	//
}

func (this *CacheBlocks) GetBlock4(host, url string) (template.HTML, bool) {
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock4(host, url string, data template.HTML) {
	//
}
