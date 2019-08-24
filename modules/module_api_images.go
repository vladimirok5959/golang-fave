package modules

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"

	"golang-fave/engine/wrapper"

	"github.com/disintegration/imaging"
)

func (this *Modules) api_GenerateImage(wrap *wrapper.Wrapper, width, height int, resize bool, filename string) ([]byte, bool, string, error) {
	file_ext := ""
	if strings.ToLower(filepath.Ext(filename)) == ".png" {
		file_ext = "image/png"
	} else if strings.ToLower(filepath.Ext(filename)) == ".jpg" {
		file_ext = "image/jpeg"
	} else if strings.ToLower(filepath.Ext(filename)) == ".jpeg" {
		file_ext = "image/jpeg"
	}

	src, err := imaging.Open(filename)
	if err != nil {
		return []byte(""), false, file_ext, err
	}

	if !resize {
		src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
	} else {
		src = imaging.Fit(src, width, height, imaging.Lanczos)
	}

	var out_bytes bytes.Buffer
	out := bufio.NewWriter(&out_bytes)

	if file_ext == "image/png" {
		imaging.Encode(out, src, imaging.PNG)
	} else if file_ext == "image/jpeg" {
		imaging.Encode(out, src, imaging.JPEG)
	} else {
		return []byte(""), false, file_ext, nil
	}

	return out_bytes.Bytes(), true, file_ext, nil
}
