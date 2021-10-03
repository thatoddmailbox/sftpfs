package sftpfs

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func parseString(data []byte) (string, int) {
	length := binary.BigEndian.Uint32(data)
	result := string(data[4 : 4+length])
	return result, int(length + 4)
}

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
	} else if p.Type == fxpName {
		r := packetFXPName{
			Entries: []packetFXPNameEntry{},
		}
		n := 4

		// special case!
		r.Count = binary.BigEndian.Uint32(p.Data[n:])
		n += 4

		for i := 0; i < int(r.Count); i++ {
			filename, sn := parseString(p.Data[n:])
			n += sn

			longname, sn := parseString(p.Data[n:])
			n += sn

			attrs, sn := parseAttrs(p.Data[n:])
			n += sn

			r.Entries = append(r.Entries, packetFXPNameEntry{
				Filename:   filename,
				Longname:   longname,
				Attributes: attrs,
			})
		}

		return r, nil
	} else if p.Type == fxpAttrs {
		r := packetFXPAttrs{}

		r.Attributes, _ = parseAttrs(p.Data[4:])

		return r, nil
	}

	return nil, fmt.Errorf("sftpfs: cannot parse packet with type %d", p.Type)
}
