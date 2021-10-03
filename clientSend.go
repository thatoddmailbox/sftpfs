package sftpfs

import (
	"encoding/binary"
	"sync/atomic"

	"golang.org/x/crypto/ssh"
)

func (c *Client) Request(packetType byte, payload interface{}) (*packet, error) {
	id := atomic.AddUint32(&c.atomicRequestID, 1)
	idBytes := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(idBytes, id)

	d := ssh.Marshal(payload)
	d = append(idBytes, d...)

	p := packet{
		Type: packetType,
		Data: d,
	}

	ch := make(chan *packet, 1)
	c.responseChannels.Store(id, ch)

	err := c.sendPacket(p)
	if err != nil {
		return nil, err
	}

	response := <-ch

	return response, nil
}

func (c *Client) sendPacket(p packet) error {
	p.Length = uint32(len(p.Data)) + 1

	_, err := c.stdin.Write(ssh.Marshal(p))
	if err != nil {
		return err
	}

	return nil
}
