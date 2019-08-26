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
	return &CacheBlocks{
		cacheBlocks: map[string]cache{},
	}
}

func (this *CacheBlocks) prepare(host string) {
	if _, ok := this.cacheBlocks[host]; !ok {
		this.cacheBlocks[host] = cache{
			CacheBlock1: map[string]template.HTML{},
			CacheBlock2: map[string]template.HTML{},
			CacheBlock3: map[string]template.HTML{},
			CacheBlock4: map[string]template.HTML{},
			CacheBlock5: map[string]template.HTML{},
		}
	}
}

func (this *CacheBlocks) Reset(host string) {
	if _, ok := this.cacheBlocks[host]; ok {
		delete(this.cacheBlocks, host)
	}
}
