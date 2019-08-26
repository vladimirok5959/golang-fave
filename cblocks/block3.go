package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) GetBlock3(host, url string) (template.HTML, bool) {
	if mapCache, ok := this.cacheBlocks[host]; ok {
		if data, ok := mapCache.CacheBlock3[url]; ok {
			return data, ok
		}
	}
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock3(host, url string, data template.HTML) {
	if _, ok := this.cacheBlocks[host]; !ok {
		this.cacheBlocks[host] = cache{}
	}
	this.cacheBlocks[host].CacheBlock3[url] = data
}
