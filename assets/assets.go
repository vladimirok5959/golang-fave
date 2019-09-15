package assets

import (
	"golang-fave/consts"

	"github.com/vladimirok5959/golang-server-resources/resource"
)

func PopulateResources(res *resource.Resource) {
	res.Add(consts.AssetsCpImgLoadGif, "image/gif", CpImgLoadGif)
	res.Add(consts.AssetsCpCodeMirrorCss, "text/css", CpCodeMirrorCss)
	res.Add(consts.AssetsCpCodeMirrorJs, "application/javascript; charset=utf-8", CpCodeMirrorJs)
	res.Add(consts.AssetsCpStylesCss, "text/css", CpStylesCss)
	res.Add(consts.AssetsCpWysiwygPellCss, "text/css", CpWysiwygPellCss)
	res.Add(consts.AssetsCpWysiwygPellJs, "application/javascript; charset=utf-8", CpWysiwygPellJs)
	res.Add(consts.AssetsLightGalleryCss, "text/css", LightGalleryCss)
	res.Add(consts.AssetsLightGalleryJs, "application/javascript; charset=utf-8", LightGalleryJs)
	res.Add(consts.AssetsSysBgPng, "image/png", SysBgPng)
	res.Add(consts.AssetsSysFaveIco, "image/x-icon", SysFaveIco)
	res.Add(consts.AssetsSysLogoPng, "image/png", SysLogoPng)
	res.Add(consts.AssetsSysLogoSvg, "image/svg+xml", SysLogoSvg)
	res.Add(consts.AssetsSysStylesCss, "text/css", SysStylesCss)
	res.Add(consts.AssetsCpScriptsJs, "application/javascript; charset=utf-8", CpScriptsJs)
}
