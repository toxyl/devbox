package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

var (
	createMissingDevboxFiles = core.CreateMissingDevboxFiles
	fileExists               = core.FileExists
	getDevboxPath            = core.GetDevboxPath
	getDevboxConfigPath      = core.GetDevboxConfigPath
	findCommand              = core.FindCommand
	must                     = core.Must
	checkError               = core.CheckError
	forceExit                = core.ForceExit
	log                      = glog.NewLoggerSimple("command")
)
