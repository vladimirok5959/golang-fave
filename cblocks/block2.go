package cblocks

import (
	"html/template"
)

func (this *CacheBlocks) GetBlock2(host, url string) (template.HTML, bool) {
	this.mutex.Lock()
	if mapCache, ok := this.cacheBlocks[host]; ok {
		if data, ok := mapCache.CacheBlock2[url]; ok {
			this.mutex.Unlock()
			return data, ok
		}
	}
	this.mutex.Unlock()
	return template.HTML(""), false
}

func (this *CacheBlocks) SetBlock2(host, url string, data template.HTML) {
	this.mutex.Lock()
	this.prepare(host)
	this.cacheBlocks[host].CacheBlock2[url] = data
	this.mutex.Unlock()
}
