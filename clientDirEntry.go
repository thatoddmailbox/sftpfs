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
	i, _ := d.Info()
	return i.Mode()
}

func (d *DirEntry) Info() (fs.FileInfo, error) {
	return &FileInfo{
		name:  d.d.Filename,
		attrs: d.d.Attributes,
	}, nil
}
