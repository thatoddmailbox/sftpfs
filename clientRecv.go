package sftpfs

import (
	"encoding/binary"
	"errors"
	"fmt"
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

func (c *Client) receiveLoop() {
	for {
		p, err := c.readPacket()
		if err != nil {
			if err == io.EOF && c.closed {
				break
			}

			panic(err)
		}

		requestID := binary.BigEndian.Uint32(p.Data)

		value, found := c.responseChannels.LoadAndDelete(requestID)
		if !found {
			panic(fmt.Errorf("sftpfs: got packet for unknown request id %d", requestID))
		}

		ch := value.(chan *packet)
		ch <- p
	}
}
