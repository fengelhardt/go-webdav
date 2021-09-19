package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"flag"
	"log"
)

var gConfig Config
var gUserConfig UserConfig
var gLogger *log.Logger

const (
	gUserDirPerm = 0775
)

var gConfFile = flag.String("c", "", "location of the config file")
