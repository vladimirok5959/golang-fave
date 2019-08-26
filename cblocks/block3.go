package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) ResetBlock3(host string) {
	//
}

func (this *CacheBlocks) GetBlock3(host, url string) (template.HTML, bool) {
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock3(host, url string, data template.HTML) {
	//
}
