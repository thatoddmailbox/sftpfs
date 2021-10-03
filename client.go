package sftpfs

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// A Client is an SFTP client that implements fs.FS.
type Client struct {
	cl *ssh.Client
	s  *ssh.Session

	stdin  io.WriteCloser
	stdout io.Reader
}

func (c *Client) Close() error {
	return c.cl.Close()
}

func (c *Client) sendPacket(p packet) error {
	p.Length = uint32(len(p.Data)) + 1

	_, err := c.stdin.Write(ssh.Marshal(p))
	if err != nil {
		return err
	}

	return nil
}

// Dial starts an SFTP connection with the given parameters.
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	var err error

	c := Client{}

	c.cl, err = ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}

	c.s, err = c.cl.NewSession()
	if err != nil {
		return nil, err
	}

	c.stdin, err = c.s.StdinPipe()
	if err != nil {
		return nil, err
	}

	c.stdout, err = c.s.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = c.s.RequestSubsystem("sftp")
	if err != nil {
		return nil, err
	}

	// sftp v3 section 4

	// send init packet
	err = c.sendPacket(packet{
		Type: fxpInit,
		Data: []byte{
			0, 0, 0, 3, // version 3
		},
	})
	if err != nil {
		return nil, err
	}

	p, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	versionPacket := packetFXPVersion{}
	err = ssh.Unmarshal(p.Data, &versionPacket)
	if err != nil {
		return nil, err
	}

	if versionPacket.Version != 3 {
		return nil, fmt.Errorf("sftpfs: unsupported sftp version %d", versionPacket.Version)
	}

	return &c, nil
}
