package fetdata

import (
	"golang-fave/utils"
)

func (this *FERData) RequestURI() string {
	return this.wrap.R.RequestURI
}

func (this *FERData) RequestURL() string {
	return this.wrap.R.URL.Path
}

func (this *FERData) RequestGET() string {
	return utils.ExtractGetParams(this.wrap.R.RequestURI)
}
