package oss

import (
	"github.com/raythorn/falcon/context"
	"github.com/raythorn/falcon/log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	// "strings"
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

	from, to := contentRange(ctx)
	if from == 0 && to == 0 {
		ctx.WriteHeader(416)
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

func contentLength(ctx *context.Context) int64 {
	lengthstr := ctx.Get("Content-Length")
	if len(lengthstr) == 0 {
		return 0
	}

	length, err := strconv.ParseInt(lengthstr, 10, 64)
	if err != nil {
		return 0
	}

	return length
}

func contentRange(ctx *context.Context) (int64, int64) {
	length := contentLength(ctx)
	rangestr := ctx.Get("Content-Range")
	if len(rangestr) == 0 {
		return 0, length
	}

	pattern := `(?P<from>\d*)-(?P<to>\d*)/(?P<total>\d*)`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(rangestr)

	var total int64 = 0
	var err error = nil

	for i, name := range reg.SubexpNames() {
		err = nil
		if name == "from" {
			from, err = strconv.ParseInt(matches[i], 10, 64)
		} else if name == "to" {
			to, err = strconv.ParseInt(matches[i], 10, 64)
		} else if name == "total" {
			total, err = strconv.ParseInt(matches[i], 10, 64)
		}

		if err != nil {
			return 0, 0
		}
	}

	if total != (to-from+1) || to >= length {
		return 0, 0
	}

	return from, to
}
