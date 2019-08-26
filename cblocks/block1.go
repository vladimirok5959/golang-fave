package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) GetBlock1(host, url string) (template.HTML, bool) {
	this.Lock()
	defer this.Unlock()
	if mapCache, ok := this.cacheBlocks[host]; ok {
		if data, ok := mapCache.CacheBlock1[url]; ok {
			return data, ok
		}
	}
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock1(host, url string, data template.HTML) {
	this.Lock()
	defer this.Unlock()
	this.prepare(host)
	this.cacheBlocks[host].CacheBlock1[url] = data
}
