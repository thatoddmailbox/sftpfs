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
		c:      c,
		path:   name,
		handle: handle.Handle,
	}, nil
}

func (c *Client) ReadDir(name string) ([]fs.DirEntry, error) {
	result := []fs.DirEntry{}

	openResponse, err := c.Request(fxpOpendir, packetFXPOpendir{
		Path: name,
	})
	if err != nil {
		return nil, err
	}

	status, ok := openResponse.(packetFXPStatus)
	if ok {
		// failure
		return nil, errors.New(status.Message)
	}

	handle, ok := openResponse.(packetFXPHandle)
	if !ok {
		// unknown
		return nil, errors.New("sftpfs: unexpected packet type")
	}

	// ok, we got a handle!
	// time to read
	for {
		readResponse, err := c.Request(fxpReaddir, packetFXPReaddir{
			Handle: handle.Handle,
		})
		if err != nil {
			return nil, err
		}

		status, ok := readResponse.(packetFXPStatus)
		if ok {
			// failure
			if status.StatusCode == fxEOF {
				// we're done
				break
			}

			return nil, errors.New(status.Message)
		}

		name, ok := readResponse.(packetFXPName)
		if !ok {
			// unknown
			return nil, errors.New("sftpfs: unexpected packet type")
		}

		for _, entry := range name.Entries {
			if entry.Filename == "." || entry.Filename == ".." {
				continue
			}

			result = append(result, &DirEntry{
				d: entry,
			})
		}
	}

	closeResponse, err := c.Request(fxpClose, packetFXPClose{
		Handle: handle.Handle,
	})
	if err != nil {
		return nil, err
	}

	status, ok = closeResponse.(packetFXPStatus)
	if !ok {
		// unknown
		return nil, errors.New("sftpfs: unexpected packet type")
	}

	if status.StatusCode != fxOK {
		return nil, errors.New(status.Message)
	}

	return result, nil
}
