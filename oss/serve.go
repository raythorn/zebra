package oss

import (
	"github.com/raythorn/falcon/context"
	// "net/http"
	// "os"
	// "path"
	"github.com/raythorn/falcon/log"
)

func ServeContent(ctx *context.Context) {

	log.Info("Serve Content")
	log.Info("Resource Name: %s", ctx.Get("resource"))
	// filepath := path.Clean(ctx.Get(OssRootKey) + "/" + ctx.Get(OssPathKey))

	// _, err := os.Stat(fp)
	// if os.IsNotExist(err) {
	// 	ctx.NotFound()
	// 	return
	// }

	// if err != nil {
	// 	ctx.WriteHeader(500)
	// 	ctx.WriteString(err.Error())
	// 	return
	// }
}
