package sftpfs

import (
	"errors"
	"io/fs"
	"path"
)

type File struct {
	c      *Client
	path   string
	handle []byte
}

func (f *File) Stat() (fs.FileInfo, error) {
	r, err := f.c.Request(fxpStat, packetFXPStat{
		Path: f.path,
	})
	if err != nil {
		return nil, err
	}

	status, ok := r.(packetFXPStatus)
	if ok {
		// failure
		return nil, errors.New(status.Message)
	}

	attrs, ok := r.(packetFXPAttrs)
	if !ok {
		// unknown
		return nil, errors.New("sftpfs: unexpected packet type")
	}

	return &FileInfo{
		name:  path.Base(f.path),
		attrs: attrs.Attributes,
	}, nil
}

func (f *File) Read(buf []byte) (int, error) {
	return 0, nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) ReadDir(n int) ([]fs.DirEntry, error) {
	if n != 0 {
		return nil, errors.New("sftpfs: readdir with n != 0 not implemented")
	}

	return f.c.ReadDir(f.path)
}
