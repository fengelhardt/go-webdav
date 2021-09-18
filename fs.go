package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

type UserDirFileSystem struct {
	RO             bool
	SingleUserMode bool
}

func newUserDirFileSystem(ro, singleUserMode bool) UserDirFileSystem {
	return UserDirFileSystem{RO: ro, SingleUserMode: singleUserMode}
}

func (udfs UserDirFileSystem) resolveVirtualPath(uri, user string) (string, error) {
	realpath := pathClean(uri)
	if !strings.HasPrefix(realpath, gConfig.URIPrefix) {
		return uri, fmt.Errorf("Request URI %q has wrong prefix", uri)
	}
	realpath = strings.TrimPrefix(realpath, gConfig.URIPrefix)
	if udfs.SingleUserMode {
		realpath = gConfig.BaseDir + "/" + realpath
	} else {
		userdir := fmt.Sprintf("%s/%s", gConfig.BaseDir, user)
		if _, err := os.Lstat(userdir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if err := os.Mkdir(userdir, gUserDirPerm); err != nil {
					return uri, fmt.Errorf("Cannot create user dir %q: %v", userdir, err)
				}
			} else {
				return uri, fmt.Errorf("Cannot virtualize user dir %q: %v", userdir, err)
			}
		}
		realpath = fmt.Sprintf("%s/%s", userdir, realpath)
	}
	return filepath.FromSlash(realpath), nil
}

func (udfs UserDirFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	//fmt.Println("Mkdir", name)
	if udfs.RO {
		return fs.ErrPermission
	}
	user := ctx.Value("user").(string)
	if realpath, err := udfs.resolveVirtualPath(name, user); err == nil {
		return os.Mkdir(realpath, perm)
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (udfs UserDirFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	//fmt.Println("OpenFile", name)
	if udfs.RO && (flag&os.O_RDWR > 0 || flag&os.O_CREATE > 0) {
		return nil, fs.ErrPermission
	}
	user := ctx.Value("user").(string)
	if realpath, err := udfs.resolveVirtualPath(name, user); err == nil {
		f, err := os.OpenFile(realpath, flag, perm)
		return f, err
	} else {
		return nil, fmt.Errorf("Real path not found: %s", err)
	}
}

func (udfs UserDirFileSystem) RemoveAll(ctx context.Context, name string) error {
	//fmt.Println("RemoveAll", name)
	if udfs.RO {
		return fs.ErrPermission
	}
	user := ctx.Value("user").(string)
	if realpath, err := udfs.resolveVirtualPath(name, user); err == nil {
		return os.RemoveAll(realpath)
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (udfs UserDirFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	//fmt.Println("Rename", oldName)
	if udfs.RO {
		return fs.ErrPermission
	}
	user := ctx.Value("user").(string)
	if realpathOld, err := udfs.resolveVirtualPath(oldName, user); err == nil {
		if realpathNew, err := udfs.resolveVirtualPath(newName, user); err == nil {
			return os.Rename(realpathOld, realpathNew)
		} else {
			return fmt.Errorf("Real path not found: %s", err)
		}
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (udfs UserDirFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	//fmt.Println("Stat", name)
	user := ctx.Value("user").(string)
	if realpath, err := udfs.resolveVirtualPath(name, user); err == nil {
		return os.Lstat(realpath)
	} else {
		info, _ := os.Lstat("")
		// TODO: remove w flags from info if in RO mode
		// we need to create a new FileInfo type for this
		return info, fmt.Errorf("Real path not found: %s", err)
	}
}
