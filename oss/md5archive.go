package oss

import (
	"github.com/raythorn/falcon/context"
)

type MD5Archive struct {
}

func (md5 *MD5Archive) Path(oss *Oss, ctx *context.Context) string {
	return ""
}
