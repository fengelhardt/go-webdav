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
}

func newUserDirFileSystem() UserDirFileSystem {
	return UserDirFileSystem{}
}

func resolveVirtualPath(uri, user string) (string, error) {
	realpath := pathClean(uri)
	if !strings.HasPrefix(realpath, gConfig.URIPrefix) {
		return uri, fmt.Errorf("Request URI %q has wrong prefix", uri)
	}
	realpath = strings.TrimPrefix(realpath, gConfig.URIPrefix)
	userdir := fmt.Sprintf("%s/%s", gConfig.BaseDir, user)
	if _, err := os.Lstat(userdir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.Mkdir(userdir, gUserDirPerm); err == nil {
				return uri, fmt.Errorf("Cannot create user dir %q: %v", userdir, err)
			}
		} else {
			return uri, fmt.Errorf("Cannot virtualize user dir %q: %v", userdir, err)
		}
	}
	realpath = fmt.Sprintf("%s/%s", userdir, realpath)
	return filepath.FromSlash(realpath), nil
}

func (fs UserDirFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	//fmt.Println("Mkdir", name)
	user := ctx.Value("user").(string)
	if realpath, err := resolveVirtualPath(name, user); err == nil {
		return os.Mkdir(realpath, perm)
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (fs UserDirFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	//fmt.Println("OpenFile", name)
	user := ctx.Value("user").(string)
	if realpath, err := resolveVirtualPath(name, user); err == nil {
		f, err := os.OpenFile(realpath, flag, perm)
		return f, err
	} else {
		return nil, fmt.Errorf("Real path not found: %s", err)
	}
}

func (fs UserDirFileSystem) RemoveAll(ctx context.Context, name string) error {
	//fmt.Println("RemoveAll", name)
	user := ctx.Value("user").(string)
	if realpath, err := resolveVirtualPath(name, user); err == nil {
		return os.RemoveAll(realpath)
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (fs UserDirFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	//fmt.Println("Rename", oldName)
	user := ctx.Value("user").(string)
	if realpathOld, err := resolveVirtualPath(oldName, user); err == nil {
		if realpathNew, err := resolveVirtualPath(newName, user); err == nil {
			return os.Rename(realpathOld, realpathNew)
		} else {
			return fmt.Errorf("Real path not found: %s", err)
		}
	} else {
		return fmt.Errorf("Real path not found: %s", err)
	}
}

func (fs UserDirFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	//fmt.Println("Stat", name)
	user := ctx.Value("user").(string)
	if realpath, err := resolveVirtualPath(name, user); err == nil {
		return os.Lstat(realpath)
	} else {
		info, _ := os.Lstat("")
		return info, fmt.Errorf("Real path not found: %s", err)
	}
}
