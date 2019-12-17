package assets

import (
	"golang-fave/consts"

	"github.com/vladimirok5959/golang-server-resources/resource"
)

func PopulateResources(res *resource.Resource) {
	res.Add(consts.AssetsCpImgLoadGif, "image/gif", CpImgLoadGif, 1)
	res.Add(consts.AssetsCpCodeMirrorCss, "text/css", CpCodeMirrorCss, 1)
	res.Add(consts.AssetsCpCodeMirrorJs, "application/javascript; charset=utf-8", CpCodeMirrorJs, 1)
	res.Add(consts.AssetsCpStylesCss, "text/css", CpStylesCss, 1)
	res.Add(consts.AssetsCpWysiwygPellCss, "text/css", CpWysiwygPellCss, 1)
	res.Add(consts.AssetsCpWysiwygPellJs, "application/javascript; charset=utf-8", CpWysiwygPellJs, 1)
	res.Add(consts.AssetsLightGalleryCss, "text/css", LightGalleryCss, 1)
	res.Add(consts.AssetsLightGalleryJs, "application/javascript; charset=utf-8", LightGalleryJs, 1)
	res.Add(consts.AssetsSysBgPng, "image/png", SysBgPng, 1)
	res.Add(consts.AssetsSysFaveIco, "image/x-icon", SysFaveIco, 1)
	res.Add(consts.AssetsSysLogoPng, "image/png", SysLogoPng, 1)
	res.Add(consts.AssetsSysLogoSvg, "image/svg+xml", SysLogoSvg, 1)
	res.Add(consts.AssetsSysStylesCss, "text/css", SysStylesCss, 1)
	res.Add(consts.AssetsCpScriptsJs, "application/javascript; charset=utf-8", CpScriptsJs, 1)
	res.Add(consts.AssetsSysPlaceholderPng, "image/png", SysPlaceholderPng, 1)
}
