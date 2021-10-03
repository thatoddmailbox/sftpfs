package sftpfs

import "encoding/binary"

const (
	attrSize        = 0x00000001
	attrUidgid      = 0x00000002
	attrPermissions = 0x00000004
	attrAcmodtime   = 0x00000008
	attrExtended    = 0x80000000
)

type attrs struct {
	Flags       uint32
	Size        uint64
	UID         uint32
	GID         uint32
	Permissions uint32
	Atime       uint32
	Mtime       uint32
	Extended    [][]string
}

func parseAttrs(data []byte) (attrs, int) {
	// see sftp v3 section 5

	n := 4
	f := binary.BigEndian.Uint32(data)
	a := attrs{
		Flags:    f,
		Extended: [][]string{},
	}

	if f&attrSize != 0 {
		a.Size = binary.BigEndian.Uint64(data[n:])
		n += 8
	}
	if f&attrUidgid != 0 {
		a.UID = binary.BigEndian.Uint32(data[n:])
		n += 4
		a.GID = binary.BigEndian.Uint32(data[n:])
		n += 4
	}
	if f&attrPermissions != 0 {
		a.Permissions = binary.BigEndian.Uint32(data[n:])
		n += 4
	}
	if f&attrAcmodtime != 0 {
		a.Atime = binary.BigEndian.Uint32(data[n:])
		n += 4
		a.Mtime = binary.BigEndian.Uint32(data[n:])
		n += 4
	}
	if f&attrExtended != 0 {
		extendedCount := binary.BigEndian.Uint32(data[n:])
		for i := 0; i < int(extendedCount); i++ {
			extendedType, sn := parseString(data[n:])
			n += sn
			extendedData, sn := parseString(data[n:])
			n += sn

			a.Extended = append(a.Extended, []string{
				extendedType, extendedData,
			})
		}
	}

	return a, n
}
