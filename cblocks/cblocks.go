package cblocks

import (
	"html/template"
)

type cache struct {
	CacheBlock1 map[string]template.HTML
	CacheBlock2 map[string]template.HTML
	CacheBlock3 map[string]template.HTML
	CacheBlock4 map[string]template.HTML
	CacheBlock5 map[string]template.HTML
}

type CacheBlocks struct {
	cacheBlocks map[string]cache
}

func New() *CacheBlocks {
	return &CacheBlocks{}
}

func (this *CacheBlocks) ResetBlock1(host string) {
	//
}

func (this *CacheBlocks) GetBlock1(host, url string) (template.HTML, bool) {
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock1(host, url string, data template.HTML) {
	//
}
