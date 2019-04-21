## Overview
file-sync is a real-time watch file changes and sync files to remote server.

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported. [See the install instructions for Go.](https://golang.org/doc/install)

```
$ go get github.com/xigang/file-sync
$ go build -o file-sync
```
Move file-sync binary to your `PATH`

## Example

```
$ file-sync -username=<username> 			\
			-ssh-key=<ssh_host_rsa_key> 	\
			-hostname=<remote-hostname>		 \
			-filename=<watch-local-filename> \
			-remote-path=<remote-server-filename>
```

## License

- [Apache License](LICENSE)
