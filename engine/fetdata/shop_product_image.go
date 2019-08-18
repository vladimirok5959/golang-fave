package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopProductImage struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_product_image
}

func (this *ShopProductImage) ProductId() int {
	if this == nil {
		return 0
	}
	return this.object.A_product_id
}

func (this *ShopProductImage) FileName() string {
	if this == nil {
		return ""
	}
	return this.object.A_filename
}

func (this *ShopProductImage) FullImage() string {
	if this == nil {
		return ""
	}
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail0() string {
	if this == nil {
		return ""
	}
	return "/api/product-image/thumb-0/" + utils.IntToStr(this.object.A_product_id) + "/" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail1() string {
	if this == nil {
		return ""
	}
	return "/api/product-image/thumb-1/" + utils.IntToStr(this.object.A_product_id) + "/" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail2() string {
	if this == nil {
		return ""
	}
	return "/api/product-image/thumb-2/" + utils.IntToStr(this.object.A_product_id) + "/" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail3() string {
	if this == nil {
		return ""
	}
	return "/api/product-image/thumb-3/" + utils.IntToStr(this.object.A_product_id) + "/" + this.object.A_filename
}

func (this *ShopProductImage) ThumbnailSize0() [2]int {
	return (*this.wrap.Config).Shop.Thumbnails.Thumbnail0
}

func (this *ShopProductImage) ThumbnailSize1() [2]int {
	return (*this.wrap.Config).Shop.Thumbnails.Thumbnail1
}

func (this *ShopProductImage) ThumbnailSize2() [2]int {
	return (*this.wrap.Config).Shop.Thumbnails.Thumbnail2
}

func (this *ShopProductImage) ThumbnailSize3() [2]int {
	return (*this.wrap.Config).Shop.Thumbnails.Thumbnail3
}
