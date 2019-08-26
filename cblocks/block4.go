package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) GetBlock4(host, url string) (template.HTML, bool) {
	this.mutex.Lock()
	if mapCache, ok := this.cacheBlocks[host]; ok {
		if data, ok := mapCache.CacheBlock4[url]; ok {
			this.mutex.Unlock()
			return data, ok
		}
	}
	this.mutex.Unlock()
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock4(host, url string, data template.HTML) {
	this.mutex.Lock()
	this.prepare(host)
	this.cacheBlocks[host].CacheBlock4[url] = data
	this.mutex.Unlock()
}
