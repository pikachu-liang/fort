package fileutil

import (
	"fmt"
	"github.com/golang/glog"
	"os"
	"path/filepath"
)

const (
	PrivateFileMode = 0600
	PrivateDirMode  = 0700
)

// IsDirWriteable checks if dir is writable by trying to write and remove a file to it.
func IsDirWriteable(dir string) error {
	f, err := filepath.Abs(filepath.Join(dir, ".tmp"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(f, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// TouchDirAll creates a directories with 0700 permission if the directory
// does not exist. TouchDirAll also ensures the given directory is writable.
func TouchDirAll(dir string) error {
	// If path is already a directory, MkdirAll does nothing and returns nil, so,
	// first check if dir exist with an expected permission mode.
	if Exist(dir) {
		err := CheckDirPermission(dir, PrivateDirMode)
		if err != nil {
			glog.Error("check file failed", err)
		}
	} else {
		err := os.MkdirAll(dir, PrivateDirMode)
		if err != nil {
			// if mkdirAll("a/text") and "text" is not
			// a directory, this will return syscall.ENOTDIR
			return err
		}
	}
	return IsDirWriteable(dir)
}

// Exist returns true if a file or directory exists.
func Exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// CheckDirPermission checks permission of an existing dir.
// Returns error if dir is empty or dir permission is not as specified.
func CheckDirPermission(dir string, perm os.FileMode) error {
	if !Exist(dir) {
		return fmt.Errorf("directory %q is empty, cannot check permission", dir)
	}
	// check the permission of the directory
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	dirPerm := dirInfo.Mode().Perm()
	if dirPerm != perm {
		err = fmt.Errorf("directory %q exist, but the permission is %q, expected permission: %q",
			dir, dirInfo.Mode(), os.FileMode(PrivateDirMode))
		return err
	}
	return nil
}
