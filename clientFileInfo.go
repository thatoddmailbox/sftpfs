package sftpfs

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	name  string
	attrs attrs
}

func (i *FileInfo) Name() string {
	return i.name
}

func (i *FileInfo) Size() int64 {
	return int64(i.attrs.Size)
}

func (i *FileInfo) Mode() fs.FileMode {
	// TODO: implement other things
	mode := fs.FileMode(0)
	if i.attrs.Permissions&(1<<14) != 0 {
		mode |= fs.ModeDir
	}
	return mode
}

func (i *FileInfo) ModTime() time.Time {
	return time.Unix(int64(i.attrs.Mtime), 0)
}

func (i *FileInfo) IsDir() bool {
	return i.Mode().IsDir()
}

func (i *FileInfo) Sys() interface{} {
	return nil
}
