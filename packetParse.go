package sftpfs

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (p *packet) Parse() (interface{}, error) {
	// TODO: this is kinda janky
	if p.Type == fxpStatus {
		r := packetFXPStatus{}

		err := ssh.Unmarshal(p.Data[4:], &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else if p.Type == fxpHandle {
		r := packetFXPHandle{}

		err := ssh.Unmarshal(p.Data[4:], &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	}

	return nil, fmt.Errorf("sftpfs: cannot parse packet with type %d", p.Type)
}
