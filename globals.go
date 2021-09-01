package main

import (
	"flag"
)

var gConfig Config
var gUserConfig UserConfig

const (
	gUserDirPerm = 0775
)

var gConfFile = flag.String("c", "", "location of the config file")
