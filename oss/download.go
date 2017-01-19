package oss

import (
	"github.com/raythorn/falcon/context"
	"net/http"
	"os"
	"path"
)

func Download(ctx *context.Context) {

	filepath := path.Clean(ctx.Get(OssRootKey) + "/" + ctx.Get(OssPathKey))

	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		ctx.NotFound()
		return
	}

	if err != nil {
		ctx.WriteHeader(500)
		ctx.WriteString(err.Error())
		return
	}
}
