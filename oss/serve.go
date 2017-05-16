package oss

import (
	"github.com/raythorn/falcon/context"
	// "net/http"
	"os"
	// "path"
	// "github.com/raythorn/falcon/log"
)

func ServeContent(ctx *context.Context) {

	respath := ctx.Get(OssPathKey)
	if len(respath) == 0 {
		resp := map[string]interface{}{}
		resp["code"] = 1
		resp["msg"] = "Invalid resid, not MD5 string"
		ctx.JSON(resp, false)
		return
	}

	if ctx.Method() == "GET" {
		if !isExist(respath) {
			ctx.NotFound()
			return
		}
	} else if ctx.Method() == "POST" {

	}
}

func isExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return false
}
