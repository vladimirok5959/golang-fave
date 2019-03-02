package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) RequestURI() string {
	return this.Wrap.R.RequestURI
}

func (this *FERData) RequestURL() string {
	return this.Wrap.R.URL.Path
}

func (this *FERData) RequestGET() string {
	return utils.ExtractGetParams(this.Wrap.R.RequestURI)
}
