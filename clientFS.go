package sftpfs

import (
	"errors"
	"io/fs"
)

func (c *Client) Open(name string) (fs.File, error) {
	r, err := c.Request(fxpOpen, packetFXPOpen{
		Filename: name,
		PFlags:   fxfRead,
		Attributes: []byte{
			0, 0, 0, 0, //flags
		},
	})
	if err != nil {
		return nil, err
	}

	status, ok := r.(packetFXPStatus)
	if ok {
		// failure
		return nil, errors.New(status.Message)
	}

	handle, ok := r.(packetFXPHandle)
	if !ok {
		// unknown
		return nil, errors.New("sftpfs: unexpected packet type")
	}

	return &File{
		handle: handle.Handle,
	}, nil
}
