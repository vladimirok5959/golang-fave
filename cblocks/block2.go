package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) ResetBlock2(host string) {
	//
}

func (this *CacheBlocks) GetBlock2(host, url string) (template.HTML, bool) {
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock2(host, url string, data template.HTML) {
	//
}
