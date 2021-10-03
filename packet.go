package sftpfs

type packet struct {
	Length uint32
	Type   byte
	Data   []byte `ssh:"rest"`
}

// sftp v3 page 6
const (
	fxpInit          = 1
	fxpVersion       = 2
	fxpOpen          = 3
	fxpClose         = 4
	fxpRead          = 5
	fxpWrite         = 6
	fxpLstat         = 7
	fxpFstat         = 8
	fxpSetstat       = 9
	fxpFsetstat      = 10
	fxpOpendir       = 11
	fxpReaddir       = 12
	fxpRemove        = 13
	fxpMkdir         = 14
	fxpRmdir         = 15
	fxpRealpath      = 16
	fxpStat          = 17
	fxpRename        = 18
	fxpReadlink      = 19
	fxpSymlink       = 20
	fxpStatus        = 101
	fxpHandle        = 102
	fxpData          = 103
	fxpName          = 104
	fxpAttrs         = 105
	fxpExtended      = 200
	fxpExtendedReply = 201
)

const (
	fxfRead   = (1 << 0)
	fxfWrite  = (1 << 1)
	fxfAppend = (1 << 2)
	fxfCreat  = (1 << 3)
	fxfTrunc  = (1 << 4)
	fxfExcl   = (1 << 5)
)

type packetFXPVersion struct {
	Version       uint32
	ExtensionData []byte `ssh:"rest"`
}

// everything below this line are requests
// they must have the ID field prepended

type packetFXPOpen struct {
	Filename   string
	PFlags     uint32
	Attributes []byte `ssh:"rest"`
}

type packetFXPStatus struct {
	StatusCode uint32
	Message    string
	Language   string
}

type packetFXPHandle struct {
	HandleLength uint32
	Handle       []byte `ssh:"rest"`
}
