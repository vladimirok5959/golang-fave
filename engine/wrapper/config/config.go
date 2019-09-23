package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Engine struct {
		MainModule int
	}
	Blog struct {
		Pagination struct {
			Index    int
			Category int
		}
	}
	Shop struct {
		Pagination struct {
			Index    int
			Category int
		}
		Thumbnails struct {
			Thumbnail0    [3]int
			Thumbnail1    [3]int
			Thumbnail2    [3]int
			Thumbnail3    [3]int
			ThumbnailFull [3]int
		}
	}
	API struct {
		XML struct {
			Enabled int
			Name    string
			Company string
			Url     string
		}
	}
}

func ConfigNew() *Config {
	c := &Config{}
	c.configDefault()
	return c
}

func (this *Config) configDefault() {
	this.Engine.MainModule = 0

	this.Blog.Pagination.Index = 5
	this.Blog.Pagination.Category = 5

	this.Shop.Pagination.Index = 9
	this.Shop.Pagination.Category = 9

	this.Shop.Thumbnails.Thumbnail0[0] = 100
	this.Shop.Thumbnails.Thumbnail0[1] = 100
	this.Shop.Thumbnails.Thumbnail0[2] = 0

	this.Shop.Thumbnails.Thumbnail1[0] = 200
	this.Shop.Thumbnails.Thumbnail1[1] = 200
	this.Shop.Thumbnails.Thumbnail1[2] = 0

	this.Shop.Thumbnails.Thumbnail2[0] = 250
	this.Shop.Thumbnails.Thumbnail2[1] = 250
	this.Shop.Thumbnails.Thumbnail2[2] = 0

	this.Shop.Thumbnails.Thumbnail3[0] = 450
	this.Shop.Thumbnails.Thumbnail3[1] = 450
	this.Shop.Thumbnails.Thumbnail3[2] = 0

	this.Shop.Thumbnails.ThumbnailFull[0] = 1000
	this.Shop.Thumbnails.ThumbnailFull[1] = 800
	this.Shop.Thumbnails.ThumbnailFull[2] = 1

	this.API.XML.Enabled = 0
	this.API.XML.Name = ""
	this.API.XML.Company = ""
	this.API.XML.Url = ""
}

func (this *Config) ConfigRead(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return dec.Decode(this)
}

func (this *Config) ConfigWrite(file string) error {
	r, err := json.Marshal(this)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(string(r))
	return err
}