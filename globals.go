package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"flag"
)

var gConfig Config
var gUserConfig UserConfig

const (
	gUserDirPerm = 0775
)

var gConfFile = flag.String("c", "", "location of the config file")
