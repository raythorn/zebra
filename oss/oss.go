package oss

import (
	"github.com/raythorn/falcon/context"
)

// OSS archive manager, which can arrange objects path with your own algrithem
type Archive interface {
	Path(ctx *context.Context) string
}

// Object storage service, handle object upload and download request
type Oss struct {
	archive *Archive
	root    string
}

func New(root string, archive *Archive) *Oss {
	oss := &Oss{root: root, archive: archive}
	return oss
}
