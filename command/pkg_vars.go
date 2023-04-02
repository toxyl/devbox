package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

var (
	createMissingDevboxFiles = core.CreateMissingDevboxFiles
	fileExists               = core.FileExists
	getStorageDir            = core.GetStorageDir
	getWorkspacePath         = core.GetWorkspacePath
	getWorkspaceConfigPath   = core.GetWorkspaceConfigPath
	getDevboxPath            = core.GetDevboxPath
	getDevboxConfigPath      = core.GetDevboxConfigPath
	findCommand              = core.FindCommand
	must                     = core.Must
	checkError               = core.CheckError
	forceExit                = core.ForceExit
	log                      = glog.NewLoggerSimple("command")
)
