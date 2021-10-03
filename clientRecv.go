package sftpfs

import (
	"encoding/binary"
	"errors"
	"io"
)

func (c *Client) readPacket() (*packet, error) {
	var buffer []byte = []byte{0, 0, 0, 0}
	_, err := io.ReadFull(c.stdout, buffer[:4])
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(buffer)
	buffer = make([]byte, length)
	_, err = io.ReadFull(c.stdout, buffer)
	if err != nil {
		return nil, err
	}

	if len(buffer) < 1 {
		return nil, errors.New("sftpfs: received packet too smal")
	}

	return &packet{
		Length: length,
		Type:   buffer[0],
		Data:   buffer[1:],
	}, nil
}
