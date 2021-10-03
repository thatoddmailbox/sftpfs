package sftpfs

import "io/fs"

type DirEntry struct {
	d packetFXPNameEntry
}

func (d *DirEntry) Name() string {
	return d.d.Filename
}

func (d *DirEntry) IsDir() bool {
	return d.Type().IsDir()
}

func (d *DirEntry) Type() fs.FileMode {
	// TODO: implement other things
	mode := fs.FileMode(0)
	if d.d.Attributes.Permissions&(1<<14) != 0 {
		mode |= fs.ModeDir
	}
	return mode
}

func (d *DirEntry) Info() (fs.FileInfo, error) {
	return &FileInfo{}, nil
}
