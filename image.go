package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang-fave/engine/wrapper/config"
	"golang-fave/utils"

	"github.com/disintegration/imaging"
)

func image_generate(width, height int, resize bool, fsrc, fdst string) {
	src, err := imaging.Open(fsrc)
	if err == nil {
		if !resize {
			src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
		} else {
			src = imaging.Fit(src, width, height, imaging.Lanczos)
		}
		if err := imaging.Save(src, fdst); err != nil {
			fmt.Printf("Image generation error: %v\n", err)
		}
	}
}

func image_create(www, src, dst, typ string, conf *config.Config) {
	width := (*conf).Shop.Thumbnails.Thumbnail0[0]
	height := (*conf).Shop.Thumbnails.Thumbnail0[1]
	resize := false

	if typ == "thumb-1" {
		width = (*conf).Shop.Thumbnails.Thumbnail1[0]
		height = (*conf).Shop.Thumbnails.Thumbnail1[1]
		if (*conf).Shop.Thumbnails.Thumbnail1[2] == 1 {
			resize = true
		}
	} else if typ == "thumb-2" {
		width = (*conf).Shop.Thumbnails.Thumbnail2[0]
		height = (*conf).Shop.Thumbnails.Thumbnail2[1]
		if (*conf).Shop.Thumbnails.Thumbnail2[2] == 1 {
			resize = true
		}
	} else if typ == "thumb-3" {
		width = (*conf).Shop.Thumbnails.Thumbnail3[0]
		height = (*conf).Shop.Thumbnails.Thumbnail3[1]
		if (*conf).Shop.Thumbnails.Thumbnail3[2] == 1 {
			resize = true
		}
	} else if typ == "thumb-full" {
		width = (*conf).Shop.Thumbnails.ThumbnailFull[0]
		height = (*conf).Shop.Thumbnails.ThumbnailFull[1]
		if (*conf).Shop.Thumbnails.ThumbnailFull[2] == 1 {
			resize = true
		}
	}

	image_generate(width, height, resize, src, dst)
}

func image_detect(www, file string, conf *config.Config) {
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
				return
			}
			if !utils.IsFileExists(file_thumb_1) {
				image_create(www, file, file_thumb_1, "thumb-1", conf)
				return
			}
			if !utils.IsFileExists(file_thumb_2) {
				image_create(www, file, file_thumb_2, "thumb-2", conf)
				return
			}
			if !utils.IsFileExists(file_thumb_3) {
				image_create(www, file, file_thumb_3, "thumb-3", conf)
				return
			}
			if !utils.IsFileExists(file_thumb_full) {
				image_create(www, file, file_thumb_full, "thumb-full", conf)
				return
			}
		}
	}
}

func image_loop(www_dir string, stop chan bool) {
	dirs, err := ioutil.ReadDir(www_dir)
	if err == nil {
		for _, dir := range dirs {
			target_dir := strings.Join([]string{www_dir, dir.Name(), "htdocs", "products", "images"}, string(os.PathSeparator))
			if utils.IsDirExists(target_dir) {
				cfile := strings.Join([]string{www_dir, dir.Name(), "config", "config.json"}, string(os.PathSeparator))
				conf := config.ConfigNew()
				if err := conf.ConfigRead(cfile); err == nil {
					pattern := target_dir + string(os.PathSeparator) + "*" + string(os.PathSeparator) + "*.*"
					if files, err := filepath.Glob(pattern); err == nil {
						for _, file := range files {
							select {
							case <-stop:
								break
							default:
								image_detect(www_dir, file, conf)
							}
						}
					}
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
		case <-time.After(8 * time.Second):
			fmt.Println("Image error: force exit by timeout after 8 seconds")
			return
		}
	}
}
