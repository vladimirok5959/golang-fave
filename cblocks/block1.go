package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) ResetBlock1(host string) {
	//
}

func (this *CacheBlocks) GetBlock1(host, url string) (template.HTML, bool) {
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock1(host, url string, data template.HTML) {
	//
}
