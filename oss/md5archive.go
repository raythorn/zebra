package oss

import (
	"github.com/raythorn/falcon/context"
	"github.com/raythorn/falcon/log"
	"path"
	"strings"
)

type MD5Archive struct {
}

func (md5 *MD5Archive) Path(oss *Oss, ctx *context.Context) string {
	category := ctx.Get("category")
	resid := ctx.Get("resid")
	root := oss.Root()

	ext := path.Ext(resid)
	id := strings.TrimSuffix(resid, ext)
	if !md5.isMd5(id) {
		log.Debug("Not Md5 String: %s", id)
		return ""
	}

	return path.Join(root, category, id[0:2], id[2:5], resid)
}

func (md5 *MD5Archive) isMd5(bytes string) bool {

	if len(bytes) != 32 {
		return false
	}

	for _, ch := range bytes {
		if (ch >= 48 && ch <= 57) || (ch >= 65 && ch <= 70) || (ch >= 97 && ch <= 102) {
			continue
		}

		return false
	}

	return true
}
