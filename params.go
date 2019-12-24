package main

import (
	"os"

	"golang-fave/engine/consts"
	"golang-fave/engine/utils"
)

func read_env_params() {
	if consts.ParamHost == "0.0.0.0" {
		if os.Getenv("FAVE_HOST") != "" {
			consts.ParamHost = os.Getenv("FAVE_HOST")
		}
	}
	if consts.ParamPort == 8080 {
		if os.Getenv("FAVE_PORT") != "" {
			consts.ParamPort = utils.StrToInt(os.Getenv("FAVE_PORT"))
		}
	}
	if consts.ParamWwwDir == "" {
		if os.Getenv("FAVE_DIR") != "" {
			consts.ParamWwwDir = os.Getenv("FAVE_DIR")
		}
	}
	if consts.ParamDebug == false {
		if os.Getenv("FAVE_DEBUG") == "true" {
			consts.ParamDebug = true
		}
	}
	if consts.ParamKeepAlive == false {
		if os.Getenv("FAVE_KEEPALIVE") == "true" {
			consts.ParamKeepAlive = true
		}
	}
}
