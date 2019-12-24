package fetdata

import (
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

type ShopProductImage struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_product_image
}

func (this *ShopProductImage) load() *ShopProductImage {
	return this
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
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/thumb-0-" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail1() string {
	if this == nil {
		return ""
	}
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/thumb-1-" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail2() string {
	if this == nil {
		return ""
	}
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/thumb-2-" + this.object.A_filename
}

func (this *ShopProductImage) Thumbnail3() string {
	if this == nil {
		return ""
	}
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/thumb-3-" + this.object.A_filename
}

func (this *ShopProductImage) ThumbnailFull() string {
	if this == nil {
		return ""
	}
	return "/products/images/" + utils.IntToStr(this.object.A_product_id) + "/thumb-full-" + this.object.A_filename
}

func (this *ShopProductImage) ThumbnailSize0() [2]int {
	return [2]int{
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail0[0],
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail0[1],
	}
}

func (this *ShopProductImage) ThumbnailSize1() [2]int {
	return [2]int{
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail1[0],
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail1[1],
	}
}

func (this *ShopProductImage) ThumbnailSize2() [2]int {
	return [2]int{
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail2[0],
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail2[1],
	}
}

func (this *ShopProductImage) ThumbnailSize3() [2]int {
	return [2]int{
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail3[0],
		(*this.wrap.Config).Shop.Thumbnails.Thumbnail3[1],
	}
}

func (this *ShopProductImage) ThumbnailSizeFull() [2]int {
	return [2]int{
		(*this.wrap.Config).Shop.Thumbnails.ThumbnailFull[0],
		(*this.wrap.Config).Shop.Thumbnails.ThumbnailFull[1],
	}
}
