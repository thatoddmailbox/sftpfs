package sftpfs

import (
	"errors"
	"io"
	"io/fs"
	"path"
)

type File struct {
	c *Client

	path   string
	handle []byte
	closed bool

	offset int64
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
	r, err := f.c.Request(fxpRead, packetFXPRead{
		Handle: f.handle,
		Offset: uint64(f.offset),
		Length: uint32(len(buf)),
	})
	if err != nil {
		return 0, err
	}

	status, ok := r.(packetFXPStatus)
	if ok {
		// failure
		if status.StatusCode == fxEOF {
			return 0, io.EOF
		}

		return 0, errors.New(status.Message)
	}

	data, ok := r.(packetFXPData)
	if !ok {
		// unknown
		return 0, errors.New("sftpfs: unexpected packet type")
	}

	f.offset += int64(len(data.Data))
	n := copy(buf, data.Data)

	return n, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekCurrent {
		f.offset += offset
		return int64(f.offset), nil
	} else if whence == io.SeekEnd {
		// TODO: can we do better?

		i, err := f.Stat()
		if err != nil {
			return int64(f.offset), err
		}

		f.offset = i.Size() + offset
		return int64(f.offset), nil
	}

	f.offset = offset
	return int64(offset), nil
}

func (f *File) Close() error {
	f.closed = true

	closeResponse, err := f.c.Request(fxpClose, packetFXPClose{
		Handle: f.handle,
	})
	if err != nil {
		return err
	}

	status, ok := closeResponse.(packetFXPStatus)
	if !ok {
		// unknown
		return errors.New("sftpfs: unexpected packet type")
	}

	if status.StatusCode != fxOK {
		return errors.New(status.Message)
	}

	return nil
}

func (f *File) ReadDir(n int) ([]fs.DirEntry, error) {
	if n > 0 {
		return nil, errors.New("sftpfs: readdir with n > 0 not implemented")
	}

	return f.c.ReadDir(f.path)
}
