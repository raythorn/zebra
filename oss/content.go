package oss

import (
	"github.com/raythorn/falcon/context"
	"github.com/raythorn/falcon/log"
	"net/http"
	"os"
	"strings"
	"time"
)

func ServeContent(ctx *context.Context) {

	respath := ctx.Get(OssPathKey)
	if len(respath) == 0 || !isExist(respath) {
		ctx.NotFound()
		return
	}

	file, err := os.Open(respath)
	if err != nil {
		ctx.NotFound()
		return
	}
	defer file.Close()

	fileinfo, err := os.Stat(respath)
	if err != nil {
		ctx.NotFound()
		return
	}

	if fileinfo.IsDir() {
		http.ServeFile(ctx.ResponseWriter(), ctx.Request(), respath)
		return
	}

	http.ServeContent(ctx.ResponseWriter(), ctx.Request(), respath, fileinfo.ModTime(), file)
}

func DepositContent(ctx *context.Context) {
	respath := ctx.Get(OssPathKey)
	if len(respath) == 0 {
		resp := map[string]interface{}{}
		resp["code"] = 1
		resp["msg"] = "Invalid resid, not MD5 string"
		ctx.JSON(resp, false)
		return
	}

	// http.Error(ctx.ResponseWriter(), "error", 500)
	expect := ctx.Get("Expect")
	if strings.Contains(expect, "100-continue") {
		log.Debug("100-continue")
		ctx.WriteHeader(100)
		time.Sleep(20 * time.Second)
	}

	ctx.WriteHeader(200)
}

func isExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return false
}
