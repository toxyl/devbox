package repo

import "github.com/toxyl/glog"

var (
	logServer = glog.NewLoggerSimple("repo-server")
	logClient = glog.NewLoggerSimple("repo-client")
)
