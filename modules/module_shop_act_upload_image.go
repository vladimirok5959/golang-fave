package modules

import (
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopUploadImage() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-upload-image",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_count := utils.Trim(wrap.R.FormValue("count"))

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_count) {
			wrap.MsgError(`Inner system error`)
			return
		}

		pf_count_int := utils.StrToInt(pf_count)
		if pf_count_int <= 0 {
			wrap.MsgError(`Inner system error`)
			return
		}

		// Proccess all files
		for i := 1; i <= pf_count_int; i++ {
			post_field_name := "file_" + utils.IntToStr(i-1)
			if file, handler, err := wrap.R.FormFile(post_field_name); err == nil {
				if handler.Filename != "" {
					if fileBytes, err := ioutil.ReadAll(file); err == nil {
						if _, _, err := image.Decode(bytes.NewReader(fileBytes)); err == nil {
							if err := os.MkdirAll(wrap.DHtdocs+string(os.PathSeparator)+"products"+string(os.PathSeparator)+"images"+string(os.PathSeparator)+pf_id, os.ModePerm); err == nil {
								target_file_name := utils.Int64ToStr(time.Now().Unix()+int64(i-1)) + filepath.Ext(handler.Filename)
								target_file_full := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + pf_id + string(os.PathSeparator) + target_file_name
								var lastID int64 = 0
								if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
									// Block rows
									if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
										return err
									}

									// Insert row
									res, err := tx.Exec(
										`INSERT INTO shop_product_images SET
											product_id = ?,
											filename = ?,
											ord = ?
										;`,
										utils.StrToInt(pf_id),
										target_file_name,
										(utils.GetCurrentUnixTimestamp() + int64(i-1)),
									)
									if err != nil {
										return err
									}

									// Get inserted post id
									lastID, err = res.LastInsertId()
									if err != nil {
										return err
									}

									// Write file to disk
									if err := ioutil.WriteFile(target_file_full, fileBytes, 0664); err != nil {
										return err
									}
									return nil
								}); err == nil {
									wrap.Write(`$('#list-images').append('<div class="attached-img" data-id="` + utils.Int64ToStr(lastID) + `"><a href="/products/images/` + pf_id + `/` + target_file_name + `" title="` + target_file_name + `" target="_blank"><img id="pimg_` + pf_id + `_` + strings.Replace(target_file_name, ".", "_", -1) + `" src="/products/images/` + pf_id + `/thumb-0-` + target_file_name + `" onerror="WaitForFave(function(){fave.ShopProductsRetryImage(this, \'pimg_` + pf_id + `_` + strings.Replace(target_file_name, ".", "_", -1) + `\');});" /></a><a class="remove" onclick="fave.ShopProductsDeleteImage(this, ` + pf_id + `, \'` + target_file_name + `\');"><svg viewBox="1 1 11 14" width="10" height="12" class="sicon" version="1.1"><path fill-rule="evenodd" d="M11 2H9c0-.55-.45-1-1-1H5c-.55 0-1 .45-1 1H2c-.55 0-1 .45-1 1v1c0 .55.45 1 1 1v9c0 .55.45 1 1 1h7c.55 0 1-.45 1-1V5c.55 0 1-.45 1-1V3c0-.55-.45-1-1-1zm-1 12H3V5h1v8h1V5h1v8h1V5h1v8h1V5h1v9zm1-10H2V3h9v1z"></path></svg></a></div>');`)
								}
							}
						}
					}
				}
				file.Close()
			}
		}

		wrap.RecreateProductImgFiles()

		wrap.RecreateProductXmlFile()

		wrap.ResetCacheBlocks()
	})
}
