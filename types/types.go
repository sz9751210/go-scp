package types

import (
	"os"
	"time"
)

type SSHhost struct {
	Alias        string
	Host         string
	Port         string
	User         string
	IdentityFile string
}

type FileDetail struct {
	Name        string
	Size        int64
	Mode        os.FileMode
	ModTime     time.Time
	IsDirectory bool
	Owner       string
	Group       string
}
