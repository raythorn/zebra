package oss

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/raythorn/falcon/context"
	"github.com/raythorn/falcon/log"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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

	//Already exist
	if isExist(respath) {
		ctx.WriteHeader(HTTP_SUCCESS)
		log.Debug("Already exist")
		return
	}

	resdir := path.Dir(respath)
	if !isExist(resdir) {
		err := os.MkdirAll(resdir, 0770)
		if err != nil {
			ctx.WriteHeader(HTTP_INTERNAL)
			return
		}
	}

	ext := path.Ext(respath)
	cachefile := strings.TrimSuffix(respath, ext) + ".cache"
	log.Debug(cachefile)

	cache, err := os.OpenFile(cachefile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		ctx.WriteHeader(HTTP_INTERNAL)
		return
	}

	var closed bool = false

	defer func() {
		if !closed {
			log.Debug("Close Cache File")
			cache.Close()
		} else {
			log.Debug("Cache File Closed")
		}
	}()

	cacheinfo, err := os.Stat(cachefile)
	if err != nil {
		ctx.WriteHeader(HTTP_INTERNAL)
		return
	}

	from, to, chunk, length := contentRange(ctx)
	if from == 0 && to == 0 {
		ctx.WriteHeader(HTTP_RANGE)
		return
	}

	cachesize := cacheinfo.Size()
	if from > cachesize {
		ctx.WriteHeader(HTTP_RANGE)
		return
	}

	data := ctx.Body()
	var datalen int64 = int64(len(data))

	log.Debug("Chunk: %d, Size: %d", chunk, datalen)
	if chunk != datalen {
		ctx.WriteHeader(HTTP_REQUEST)
		return
	}

	size, err := cache.WriteAt(data, from)
	if size != len(data) || err != nil {
		ctx.WriteHeader(HTTP_INTERNAL)
		return
	}

	if (to + 1) == length {
		md5str := md5sum(cache)
		filename := strings.TrimSuffix(path.Base(respath), ext)
		if md5str == filename {
			cache.Close()
			closed = true
			err := os.Rename(cachefile, respath)
			if err != nil {
				ctx.WriteHeader(HTTP_INTERNAL)
			} else {
				ctx.WriteHeader(HTTP_SUCCESS)
			}
			return
		}
	}

	ctx.WriteHeader(HTTP_INTERNAL)
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

func contentRange(ctx *context.Context) (int64, int64, int64, int64) {
	length := contentLength(ctx)
	rangestr := ctx.Get("Content-Range")
	if len(rangestr) == 0 {
		return 0, length - 1, length, length
	}

	pattern := `(?P<from>\d*)-(?P<to>\d*)/(?P<total>\d*)`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(rangestr)

	var (
		from  int64 = 0
		to    int64 = 0
		chunk int64 = 0
	)

	var err error = nil

	for i, name := range reg.SubexpNames() {
		err = nil
		if name == "from" {
			from, err = strconv.ParseInt(matches[i], 10, 64)
		} else if name == "to" {
			to, err = strconv.ParseInt(matches[i], 10, 64)
		} else if name == "total" {
			chunk, err = strconv.ParseInt(matches[i], 10, 64)
		}

		if err != nil {
			return 0, 0, 0, length
		}
	}

	if chunk != (to-from+1) || from >= length || to >= length {
		return 0, 0, 0, length
	}

	return from, to, chunk, length
}

func md5sum(f *os.File) string {
	offset, err := f.Seek(0, 0)
	if err != nil || offset != 0 {
		return ""
	}

	r := bufio.NewReader(f)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%02x", h.Sum(nil))
}
