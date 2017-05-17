package oss

import (
	"github.com/raythorn/falcon/context"
	"github.com/raythorn/falcon/log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

//Keys of oss elements, this will save in context which can be referred by upload/download handler
const (
	//relative path of current file
	OssPathKey = "com.raythorn.falcon.oss.path"
)

// OSS archive manager, which can arrange objects path with your own algrithem
type Archive interface {
	Path(oss *Oss, ctx *context.Context) string
}

// Object storage service, handle object upload and download request
type Oss struct {
	root    string
	archive Archive
}

func New(root string, archive Archive) *Oss {

	log.Debug(applicationPath())

	if root == "" || !path.IsAbs(root) {
		root = path.Clean(applicationPath() + "/" + root)
	}

	log.Debug(root)

	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		err = os.MkdirAll(root, 0770)
		if err != nil {
			log.Fatal("Create root directory fail")
		}
	}

	fi, err := os.Stat(root)
	if err != nil {
		log.Fatal("Cannot stat root directory")
	}

	if !fi.IsDir() || (fi.Mode()&0700) != 0700 {
		log.Fatal("Root is not a directory or does not have read/write permission")
	}

	oss := &Oss{root: root, archive: archive}
	return oss
}

func (oss *Oss) Root() string {
	return oss.root
}

func (oss *Oss) Archive() Archive {
	return oss.archive
}

func applicationPath() string {

	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal("Cannot find application path!")
	}

	fp, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		log.Fatal("Cannot find application path!")
	}

	return fp
}
