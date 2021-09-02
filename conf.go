package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Config struct {
	UsersFile   string `json:"UsersFile"`
	TLSCertFile string `json:"TLSCertFile"`
	TLSKeyFile  string `json:"TLSKeyFile"`
	Realm       string `json:"Realm"`
	BaseDir     string `json:"BaseDir"`
	URIPrefix   string `json:"URIPrefix"`
	Port        uint16 `json:"Port"`
}

func staticConfFiles() []string {
	return []string{
		"/etc/go-webdav/conf.hjson",
		"~/.go-webdav.hjson",
		"go-webdav.hjson",
	}
}

func loadConfig(file string) (Config, error) {
	var ret Config
	trials := []string{file}
	trials = append(trials, staticConfFiles()...)
	for _, f := range trials {
		data, err := os.ReadFile(f)
		if err == nil {
			if err := unmarshalHJson(data, &ret); err != nil {
				return ret, fmt.Errorf("Error with config file %q: %s", f, err)
			} else {
				err = checkConfig(&ret)
				if err != nil {
					return ret, fmt.Errorf("Error with config in %q: %s", f, err)
				}
				return ret, nil
			}
		}
	}
	return ret, fmt.Errorf("No config file found at %q", file)
}

func checkConfig(c *Config) error {
	// check URIPrefix
	c.URIPrefix = fmt.Sprintf("/%s/", c.URIPrefix)
	c.URIPrefix = pathClean(c.URIPrefix)
	// check BaseDir
	filepath.Clean(c.BaseDir)
	if info, err := os.Lstat(c.BaseDir); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("BaseDir %q does not point to a directory", c.BaseDir)
		}
		var uid, gid int
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			// we are on linux/mac, check for directory permisions
			uid = int(stat.Uid)
			gid = int(stat.Gid)
			if uid != os.Getuid() || gid != os.Getgid() {
				return fmt.Errorf("BaseDir %q must have the same owner and group as this process", c.BaseDir)
			}
			if info.Mode().Perm()&0700 != 0700 {
				return fmt.Errorf("BaseDir %q must offer rwx permissions for its owner", c.BaseDir)
			}
		}
	} else {
		return err
	}
	return nil
}
