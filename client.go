package sftpfs

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

type subsystemRequestMsg struct {
	Subsystem string
}

// A Client is an SFTP client that implements fs.FS.
type Client struct {
	cl *ssh.Client

	ch  ssh.Channel
	req <-chan *ssh.Request
}

func (c *Client) Close() error {
	return c.cl.Close()
}

// Dial starts an SFTP connection with the given parameters.
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	var err error

	c := Client{}

	c.cl, err = ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}

	c.ch, c.req, err = c.cl.OpenChannel("session", nil)
	if err != nil {
		return nil, err
	}

	result, err := c.ch.SendRequest("subsystem", true, ssh.Marshal(subsystemRequestMsg{"sftp"}))
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, errors.New("sftpfs: server failed to start sftp subsystem, is it supported?")
	}

	return &c, nil
}
