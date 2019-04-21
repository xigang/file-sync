package main

import (
	"fmt"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	username string
	password string
	sshKey   string
	hostname string
}

type FileSyncManager struct {
	filename             string
	remotePath           string
	remotePathPermission string
	sshConfig            SSHConfig
}

func NewFileSyncManager(filename string, remotePath string, remotePathPermission string, config SSHConfig) *FileSyncManager {
	return &FileSyncManager{
		filename:             filename,
		remotePath:           remotePath,
		remotePathPermission: remotePathPermission,
		sshConfig:            config,
	}
}

func (fs *FileSyncManager) SyncFile(filename string) error {
	var (
		clientConfig ssh.ClientConfig
		err          error
	)
	if fs.sshConfig.sshKey != "" {
		clientConfig, err = auth.PrivateKey(fs.sshConfig.username, fs.sshConfig.sshKey, ssh.InsecureIgnoreHostKey())
	} else if fs.sshConfig.password != "" {
		clientConfig, err = auth.PasswordKey(fs.sshConfig.username, fs.sshConfig.password, ssh.InsecureIgnoreHostKey())
	} else {
		return fmt.Errorf("params invalid")
	}

	//create a new SCP client
	remoteSSHDAddr := fmt.Sprintf("%s:22", fs.sshConfig.hostname)
	client := scp.NewClient(remoteSSHDAddr, &clientConfig)

	//connect to the remote server
	if err = client.Connect(); err != nil {
		return fmt.Errorf("Couldn't establisch a connection to the remote server ", err)
	}
	defer client.Close()

	//open sync file
	var f *os.File
	if f, err = os.Open(fs.filename); err != nil {
		return err
	}
	defer f.Close()

	//finaly, copy the file over
	if err = client.CopyFile(f, fs.remotePath, fs.remotePathPermission); err != nil {
		return err
	}

	return nil
}
