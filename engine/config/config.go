package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Engine struct {
		MainModule  int
		Maintenance int
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
		Price struct {
			Format int
			Round  int
		}
		Orders struct {
			RequiredFields struct {
				LastName     int
				FirstName    int
				MiddleName   int
				MobilePhone  int
				EmailAddress int
				Delivery     int
				Comment      int
			}
			NotifyEmail            string
			Enabled                int
			NewOrderEmailThemeCp   string
			NewOrderEmailThemeUser string
		}
		CustomFields struct {
			Field1 struct {
				Enabled int
				Caption string
			}
			Field2 struct {
				Enabled int
				Caption string
			}
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
	SMTP struct {
		Host     string
		Port     int
		Login    string
		Password string
	}
	Modules struct {
		Enabled struct {
			Blog int
			Shop int
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
	this.Engine.Maintenance = 0

	this.Blog.Pagination.Index = 5
	this.Blog.Pagination.Category = 5

	this.Shop.Pagination.Index = 9
	this.Shop.Pagination.Category = 9

	this.Shop.Thumbnails.Thumbnail0[0] = 100
	this.Shop.Thumbnails.Thumbnail0[1] = 100
	this.Shop.Thumbnails.Thumbnail0[2] = 0

	this.Shop.Thumbnails.Thumbnail1[0] = 200
	this.Shop.Thumbnails.Thumbnail1[1] = 200
	this.Shop.Thumbnails.Thumbnail1[2] = 2

	this.Shop.Thumbnails.Thumbnail2[0] = 250
	this.Shop.Thumbnails.Thumbnail2[1] = 250
	this.Shop.Thumbnails.Thumbnail2[2] = 2

	this.Shop.Thumbnails.Thumbnail3[0] = 450
	this.Shop.Thumbnails.Thumbnail3[1] = 450
	this.Shop.Thumbnails.Thumbnail3[2] = 2

	this.Shop.Thumbnails.ThumbnailFull[0] = 1000
	this.Shop.Thumbnails.ThumbnailFull[1] = 800
	this.Shop.Thumbnails.ThumbnailFull[2] = 1

	this.Shop.Price.Format = 2
	this.Shop.Price.Round = 0

	this.Shop.Orders.RequiredFields.LastName = 1
	this.Shop.Orders.RequiredFields.FirstName = 1
	this.Shop.Orders.RequiredFields.MiddleName = 0
	this.Shop.Orders.RequiredFields.MobilePhone = 0
	this.Shop.Orders.RequiredFields.EmailAddress = 1
	this.Shop.Orders.RequiredFields.Delivery = 0
	this.Shop.Orders.RequiredFields.Comment = 0

	this.Shop.Orders.NotifyEmail = ""
	this.Shop.Orders.Enabled = 1

	this.Shop.Orders.NewOrderEmailThemeCp = "❤️ New order"
	this.Shop.Orders.NewOrderEmailThemeUser = "❤️ Thanks for your order"

	this.Shop.CustomFields.Field1.Enabled = 0
	this.Shop.CustomFields.Field1.Caption = "Custom field 1"

	this.Shop.CustomFields.Field2.Enabled = 0
	this.Shop.CustomFields.Field2.Caption = "Custom field 2"

	this.API.XML.Enabled = 0
	this.API.XML.Name = ""
	this.API.XML.Company = ""
	this.API.XML.Url = ""

	this.SMTP.Host = ""
	this.SMTP.Port = 587
	this.SMTP.Login = ""
	this.SMTP.Password = ""

	this.Modules.Enabled.Blog = 1
	this.Modules.Enabled.Shop = 1
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
