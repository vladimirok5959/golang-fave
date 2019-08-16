package wrapper

import (
	"encoding/json"
	"os"
)

type Config struct {
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
			ControlPanel [2]int
			Thumbnail1   [2]int
			Thumbnail2   [2]int
			Thumbnail3   [2]int
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

func configNew() *Config {
	c := &Config{}
	c.configDefault()
	return c
}

func (this *Config) configDefault() {
	this.Blog.Pagination.Index = 5
	this.Blog.Pagination.Category = 5

	this.Shop.Pagination.Index = 9
	this.Shop.Pagination.Category = 9

	this.Shop.Thumbnails.ControlPanel[0] = 100
	this.Shop.Thumbnails.ControlPanel[1] = 100
	this.Shop.Thumbnails.Thumbnail1[0] = 200
	this.Shop.Thumbnails.Thumbnail1[1] = 200
	this.Shop.Thumbnails.Thumbnail2[0] = 250
	this.Shop.Thumbnails.Thumbnail2[1] = 250
	this.Shop.Thumbnails.Thumbnail3[0] = 450
	this.Shop.Thumbnails.Thumbnail3[1] = 450

	this.API.XML.Enabled = 0
	this.API.XML.Name = ""
	this.API.XML.Company = ""
	this.API.XML.Url = ""
}

func (this *Config) configRead(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return dec.Decode(this)
}

func (this *Config) configWrite(file string) error {
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
