package oss

import (
	"github.com/raythorn/falcon/context"
)

//Keys of oss elements, this will save in context which can be referred by upload/download handler
const (
	//root path of current oss object
	OssRootKey = "com.raythorn.falcon.oss.root"
	//relative path of current file
	OssPathKey = "com.raythorn.falcon.oss.path"
)

// OSS archive manager, which can arrange objects path with your own algrithem
type Archive interface {
	Path(ctx *context.Context) string
}

// Object storage service, handle object upload and download request
type Oss struct {
	archive Archive
	root    string
}

func New(root string, archive Archive) *Oss {
	oss := &Oss{root: root, archive: archive}
	return oss
}

func (oss *Oss) Root() string {
	return oss.root
}

func (oss *Oss) Archive() Archive {
	return oss.archive
}
