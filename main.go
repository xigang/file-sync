package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

var (
	filename             string
	username             string
	password             string
	sshKey               string
	hostname             string
	remotePath           string
	remotePathPermission string
)

func init() {
	flag.StringVar(&filename, "filename", "", "watcher file name")
	flag.StringVar(&username, "username", "", "ssh username")
	flag.StringVar(&password, "password", "", "ssh password")
	flag.StringVar(&sshKey, "ssh-key", "", "ssh key filename")
	flag.StringVar(&hostname, "hostname", "", "ssh remote hostname")
	flag.StringVar(&remotePath, "remote-path", "/", "copy the path to the remote server")
	flag.StringVar(&remotePathPermission, "remote-path-permission", "0655", "copy the path permission to the remote server")
}

func main() {
	flag.Parse()

	var err error
	if err = flagParamsValid(); err != nil {
		glog.Errorf("%+v", err)
		return
	}

	config := SSHConfig{
		username: username,
		password: password,
		sshKey:   sshKey,
		hostname: hostname,
	}
	fileSyncManager := NewFileSyncManager(filename, remotePath, remotePathPermission, config)

	var watcher *fsnotify.Watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		glog.Errorf("failed to init watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(filename)
	if err != nil {
		glog.Errorf("failed to add watch filename: %v\n", filename)
		return
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					fileSyncManager.SyncFile(event.Name)
					time.Sleep(10 * time.Second)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				glog.Errorf("failed to watcher error: %v\n", err)
			}
		}
	}()

	<-done
}

func flagParamsValid() error {
	if filename == "" {
		return fmt.Errorf("filename is empty, please specify a filename")
	}

	if username == "" {
		return fmt.Errorf("username is empty, please specify a username")
	}

	if hostname == "" {
		return fmt.Errorf("hostname is empty, please specity a remote hostname")
	}

	if password == "" && sshKey == "" {
		return fmt.Errorf("password or ssh-key is empty, please specify password or ssh-key")
	}

	return nil
}
