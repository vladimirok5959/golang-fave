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
}

func configNew() *Config {
	c := &Config{}
	c.configDefault()
	return c
}

func (this *Config) configDefault() {
	this.Blog.Pagination.Index = 5
	this.Blog.Pagination.Category = 5
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
