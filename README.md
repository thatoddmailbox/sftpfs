# sftpfs [![Build](https://github.com/thatoddmailbox/sftpfs/actions/workflows/build.yml/badge.svg)](https://github.com/thatoddmailbox/sftpfs/actions/workflows/build.yml)

An SFTP client that implements [fs.FS](https://pkg.go.dev/io/fs#FS). It relies on the [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) package for the underlying SSH implementation.

Note that this can only be used to read files, as fs.FS only supports reads. In the future this could be extended to support writes, but there currently aren't any plans for that.

## Usage
Use the `Dial` method, like so:
```
c, err := sftpfs.Dial("tcp", "some-host-with-an-sftp-server.com:22", &ssh.ClientConfig{
	User: "username",
	Auth: []ssh.AuthMethod{
		ssh.Password("passw0rd"),
	},

	// InsecureIgnoreHostKey is unsafe!
	// You should probably use ssh.FixedHostKey.
	// See the golang.org/x/crypto/ssh documentation for details.
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
	panic(err)
}
defer c.Close()
```

`c` is an `*sftpfs.Client`, which implements fs.FS. In other words, you can pass it directly to anything that uses an fs.FS, like [fsbrowse](https://github.com/thatoddmailbox/fsbrowse).