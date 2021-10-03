package sftpfs

import (
	"io/fs"
	"time"
)

type FileInfo struct {
}

func (i *FileInfo) Name() string {
	// TODO: implement
	return ""
}

func (i *FileInfo) Size() int64 {
	// TODO: implement
	return 0
}

func (i *FileInfo) Mode() fs.FileMode {
	// TODO: implement
	return 0
}

func (i *FileInfo) ModTime() time.Time {
	// TODO: implement
	return time.Now()
}

func (i *FileInfo) IsDir() bool {
	// TODO: implement
	return false
}

func (i *FileInfo) Sys() interface{} {
	return nil
}
