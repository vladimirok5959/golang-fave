package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang-fave/engine/wrapper/config"
	"golang-fave/utils"

	"github.com/disintegration/imaging"
)

func image_generate(width, height int, resize int, fsrc, fdst string) {
	src, err := imaging.Open(fsrc)
	if err == nil {
		if resize == 0 {
			src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
			if err := imaging.Save(src, fdst); err != nil {
				fmt.Printf("Image generation error (1): %v\n", err)
			}
		} else if resize == 1 {
			src = imaging.Fit(src, width, height, imaging.Lanczos)
			if err := imaging.Save(src, fdst); err != nil {
				fmt.Printf("Image generation error (2): %v\n", err)
			}
		} else {
			src = imaging.Fit(src, width, height, imaging.Lanczos)
			dst := imaging.New(width, height, color.NRGBA{255, 255, 255, 255})
			x := 0
			y := 0
			if src.Bounds().Dx() < width {
				x = int((width - src.Bounds().Dx()) / 2)
			}
			if src.Bounds().Dy() < height {
				y = int((height - src.Bounds().Dy()) / 2)
			}
			dst = imaging.Paste(dst, src, image.Pt(x, y))
			if err := imaging.Save(dst, fdst); err != nil {
				fmt.Printf("Image generation error (3): %v\n", err)
			}
			return
		}
	}
}

func image_create(www, src, dst, typ string, conf *config.Config) {
	width := (*conf).Shop.Thumbnails.Thumbnail0[0]
	height := (*conf).Shop.Thumbnails.Thumbnail0[1]
	resize := 0

	if typ == "thumb-1" {
		width = (*conf).Shop.Thumbnails.Thumbnail1[0]
		height = (*conf).Shop.Thumbnails.Thumbnail1[1]
		resize = (*conf).Shop.Thumbnails.Thumbnail1[2]
	} else if typ == "thumb-2" {
		width = (*conf).Shop.Thumbnails.Thumbnail2[0]
		height = (*conf).Shop.Thumbnails.Thumbnail2[1]
		resize = (*conf).Shop.Thumbnails.Thumbnail2[2]
	} else if typ == "thumb-3" {
		width = (*conf).Shop.Thumbnails.Thumbnail3[0]
		height = (*conf).Shop.Thumbnails.Thumbnail3[1]
		resize = (*conf).Shop.Thumbnails.Thumbnail3[2]
	} else if typ == "thumb-full" {
		width = (*conf).Shop.Thumbnails.ThumbnailFull[0]
		height = (*conf).Shop.Thumbnails.ThumbnailFull[1]
		resize = (*conf).Shop.Thumbnails.ThumbnailFull[2]
	}

	image_generate(width, height, resize, src, dst)
}

func image_detect(www, file string, conf *config.Config) bool {
	result := false
	index := strings.LastIndex(file, string(os.PathSeparator))
	if index != -1 {
		file_name := file[index+1:]
		if !strings.HasPrefix(file_name, "thumb-") {
			file_thumb_0 := file[:index+1] + "thumb-0-" + file_name
			file_thumb_1 := file[:index+1] + "thumb-1-" + file_name
			file_thumb_2 := file[:index+1] + "thumb-2-" + file_name
			file_thumb_3 := file[:index+1] + "thumb-3-" + file_name
			file_thumb_full := file[:index+1] + "thumb-full-" + file_name
			if !utils.IsFileExists(file_thumb_0) {
				image_create(www, file, file_thumb_0, "thumb-0", conf)
				result = true
			}
			if !utils.IsFileExists(file_thumb_1) {
				image_create(www, file, file_thumb_1, "thumb-1", conf)
				result = true
			}
			if !utils.IsFileExists(file_thumb_2) {
				image_create(www, file, file_thumb_2, "thumb-2", conf)
				result = true
			}
			if !utils.IsFileExists(file_thumb_3) {
				image_create(www, file, file_thumb_3, "thumb-3", conf)
				result = true
			}
			if !utils.IsFileExists(file_thumb_full) {
				image_create(www, file, file_thumb_full, "thumb-full", conf)
				result = true
			}
		}
	}
	return result
}

func image_loop(www_dir string, stop chan bool) {
	if dirs, err := ioutil.ReadDir(www_dir); err == nil {
		for _, dir := range dirs {
			trigger := strings.Join([]string{www_dir, dir.Name(), "tmp", "trigger.img.run"}, string(os.PathSeparator))
			if utils.IsFileExists(trigger) {
				processed := false
				conf := config.ConfigNew()
				if err := conf.ConfigRead(strings.Join([]string{www_dir, dir.Name(), "config", "config.json"}, string(os.PathSeparator))); err == nil {
					target_dir := strings.Join([]string{www_dir, dir.Name(), "htdocs", "products", "images"}, string(os.PathSeparator))
					if utils.IsDirExists(target_dir) {
						pattern := target_dir + string(os.PathSeparator) + "*" + string(os.PathSeparator) + "*.*"
						if files, err := filepath.Glob(pattern); err == nil {
							for _, file := range files {
								select {
								case <-stop:
									break
								default:
									if image_detect(www_dir, file, conf) {
										if !processed {
											processed = true
										}
									}
								}
							}
						}
					}
				}
				if !processed {
					os.Remove(trigger)
				}
			}
		}
	}
}

func image_start(www_dir string) (chan bool, chan bool) {
	ch := make(chan bool)
	stop := make(chan bool)

	// Run at start
	image_loop(www_dir, stop)

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				// Run every 2 seconds
				image_loop(www_dir, stop)
			case <-ch:
				ch <- true
				return
			}
		}
	}()
	return ch, stop
}

func image_stop(ch, stop chan bool) {
	for {
		select {
		case stop <- true:
		case ch <- true:
			<-ch
			return
		case <-time.After(3 * time.Second):
			fmt.Println("Image error: force exit by timeout after 3 seconds")
			return
		}
	}
}
