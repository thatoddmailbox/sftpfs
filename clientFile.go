package sftpfs

import (
	"errors"
	"io/fs"
)

type File struct {
	handle []byte
}

func (f *File) Stat() (fs.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (f *File) Read(buf []byte) (int, error) {
	return 0, nil
}

func (f *File) Close() error {
	return nil
}
